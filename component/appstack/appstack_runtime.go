// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package appstack

import (
	"context"
	"fmt"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/sirupsen/logrus"
)

var __DefaultAppStackRuntime *AppStackRuntime

func InitAppStack(re typex.RuleX) *AppStackRuntime {
	__DefaultAppStackRuntime = &AppStackRuntime{
		RuleEngine:   re,
		Applications: make(map[string]*Application),
	}
	return __DefaultAppStackRuntime
}
func AppRuntime() *AppStackRuntime {
	return __DefaultAppStackRuntime
}

/*
*
* 加载本地文件到lua虚拟机, 但是并不执行
*
 */
func LoadApp(app *Application, luaSource string) error {

	// 重新读
	app.VM().DoString(string(luaSource))
	// 检查函数入口
	AppMainVM := app.VM().GetGlobal("Main")
	if AppMainVM == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMainVM.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	// 抽取main
	fMain := *AppMainVM.(*lua.LFunction)
	app.SetMainFunc(&fMain)
	// 加载库
	// LoadAppLib(app, __DefaultAppStackRuntime.RuleEngine)
	LoadAppLibGroup(app, __DefaultAppStackRuntime.RuleEngine)
	// 加载到内存里
	__DefaultAppStackRuntime.Applications[app.UUID] = app
	return nil
}

/*
* 此时才是真正的启动入口:
* 启动 function Main(args) --do-some-thing-- return 0 end
*
 */
func StartApp(uuid string) error {
	app, ok := __DefaultAppStackRuntime.Applications[uuid]
	if !ok {
		return fmt.Errorf("Application not exists:%s", uuid)
	}
	if app.AppState == 1 {
		return fmt.Errorf("Application already started:%s", uuid)
	}
	args := lua.LBool(false) // Main的参数，未来准备扩展
	ctx, cancel := context.WithCancel(typex.GCTX)
	app.SetCnC(ctx, cancel)
	go func(app *Application) {
		defer func() {
			glogger.GLogger.Debug("Application exit:", app.UUID)
			app.VM().Pop(1) // 防止registry溢出
			app.AppState = 0
		}()
		glogger.GLogger.Debugf("Ready to run Application:%s", app.UUID)
		app.AppState = 1
		err := app.VM().CallByParam(lua.P{
			Fn:      app.GetMainFunc(),
			NRet:    1,
			Protect: true,
			Handler: &lua.LFunction{
				GFunction: func(*lua.LState) int {
					glogger.GLogger.Debug("Protect Mode Call")
					return 0
				},
			},
		}, args)
		if err == nil {
			if app.KilledBy == "RULEX" {
				glogger.GLogger.Infof("Application %s Killed By RULEX", app.UUID)
			}
			if app.KilledBy == "NORMAL" || app.KilledBy == "" {
				glogger.GLogger.Infof("Application %s NORMAL Exited", app.UUID)
			}
			return
		}
		Debugger, Ok := app.vm.GetStack(1)
		if Ok {
			LValue, _ := app.vm.GetInfo("f", Debugger, lua.LNil)
			app.vm.GetInfo("l", Debugger, lua.LNil)
			app.vm.GetInfo("S", Debugger, lua.LNil)
			app.vm.GetInfo("u", Debugger, lua.LNil)
			app.vm.GetInfo("n", Debugger, lua.LNil)
			LFunction := LValue.(*lua.LFunction)
			LastCall := lua.DbgCall{
				Name: "_main", Pc: 0,
			}
			if len(LFunction.Proto.DbgCalls) > 0 {
				LastCall = LFunction.Proto.DbgCalls[0]
			}
			glogger.GLogger.WithFields(logrus.Fields{
				"topic": "app/console/" + uuid,
			}).Warnf("Stacktrace: Current Function Name: [%s],"+
				"What(lua|native): [%s], Source Line: [%d],"+
				" Last Call: [%s], Error message: %s",
				Debugger.Name, Debugger.What, Debugger.CurrentLine,
				LastCall.Name, err.Error(),
			)
		}
		//
		// 检查是自己死的还是被RULEX杀死
		// 1 正常结束
		// 2 被rulex删除
		// 3 跑飞了

		// 中间出现异常挂了，此时要根据: auto start 来判断是否抢救
		time.Sleep(5 * time.Second)
		if app.KilledBy == "RULEX" {
			glogger.GLogger.Infof("App %s Killed By RULEX, No need to rescue", app.UUID)
			return
		}
		if app.KilledBy == "NORMAL" {
			glogger.GLogger.Infof("App %s NORMAL Exited, No need to rescue", app.UUID)
			return
		}

		glogger.GLogger.Warnf("App %s Exited With error: %s, Maybe accident, Try to survive",
			app.UUID, err.Error())
		// TODO 到底要不要设置一个尝试重启的阈值？
		// if tryTimes >= Max -> return
		if app.AutoStart {
			glogger.GLogger.Warnf("App %s Try to restart", app.UUID)
			go StartApp(uuid)
			return
		}
		glogger.GLogger.Infof("App %s not need to restart", app.UUID)

	}(app)
	glogger.GLogger.Info("App started:", app.UUID)
	return nil
}

/*
*
* 从内存里面删除APP
*
 */
func RemoveApp(uuid string) error {
	if app, ok := __DefaultAppStackRuntime.Applications[uuid]; ok {
		app.Remove()
		delete(__DefaultAppStackRuntime.Applications, uuid)
	}
	glogger.GLogger.Info("App removed:", uuid)
	return nil
}

/*
*
* 停止应用并不删除应用, 将其进程结束，状态置0
*
 */
func StopApp(uuid string) error {
	if app, ok := __DefaultAppStackRuntime.Applications[uuid]; ok {
		app.Stop()
	}
	glogger.GLogger.Info("App removed:", uuid)
	return nil
}

/*
*
* 更新应用信息
*
 */
func UpdateApp(app Application) error {
	if oldApp, ok := __DefaultAppStackRuntime.Applications[app.UUID]; ok {
		oldApp.Name = app.Name
		oldApp.AutoStart = app.AutoStart
		oldApp.Version = app.Version
		glogger.GLogger.Info("App updated:", app.UUID)
		return nil
	}
	return fmt.Errorf("update failed, app not exists:%s", app.UUID)

}
func GetApp(uuid string) *Application {
	if app, ok := __DefaultAppStackRuntime.Applications[uuid]; ok {
		return app
	}
	return nil
}

/*
*
* 获取列表
*
 */
func AppCount() int {
	return len(__DefaultAppStackRuntime.Applications)
}
func AllApp() []*Application {
	return ListApp()
}
func ListApp() []*Application {
	apps := []*Application{}
	for _, v := range __DefaultAppStackRuntime.Applications {
		apps = append(apps, v)
	}
	return apps
}

func Stop() {
	for _, app := range __DefaultAppStackRuntime.Applications {
		glogger.GLogger.Info("Stop App:", app.UUID)
		app.Stop()
		glogger.GLogger.Info("Stop App:", app.UUID, " Successfully")
	}
	glogger.GLogger.Info("Appstack stopped")

}
func GetRuleX() typex.RuleX {
	return __DefaultAppStackRuntime.RuleEngine
}
