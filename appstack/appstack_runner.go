package appstack

import (
	"context"
	"fmt"
	"os"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

/*
*
* 管理器
*
 */
type AppStack struct {
	re           typex.RuleX
	applications map[string]*typex.Application
}

func NewAppStack(re typex.RuleX) typex.XAppStack {
	as := new(AppStack)
	as.re = re
	as.applications = map[string]*typex.Application{}
	return as
}

/*
*
* 加载本地文件到lua虚拟机, 但是并不执行
*
 */
func (as *AppStack) LoadApp(app *typex.Application) error {
	bytes, err := os.ReadFile(app.Filepath)
	if err != nil {
		return err
	}
	// 重新读
	app.VM().DoString(string(bytes))
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
	LoadAppLib(app, as.re)
	// 加载到内存里
	as.applications[app.UUID] = app
	return nil
}

/*
* 此时才是真正的启动入口:
* 启动 function Main(args) --do-some-thing-- return 0 end
*
 */
func (as *AppStack) StartApp(uuid string) error {
	app, ok := as.applications[uuid]
	if !ok {
		return fmt.Errorf("app not exists:%s", uuid)
	}
	// args := lua.LBool(false) // Main的参数，未来准备扩展
	ctx, cancel := context.WithCancel(typex.GCTX)
	app.SetCnC(ctx, cancel)
	go func(ctx context.Context) {
		appId := app.UUID
		defer func() {
			app.AppState = 0
			glogger.GLogger.Debug("App exit:", appId)
		}()
		app.VM().SetContext(ctx)
		glogger.GLogger.Debug("Ready to run app:", app.UUID, app.Name, app.Version)
		app.AppState = 1
		err := app.VM().CallByParam(lua.P{
			Fn:      app.GetMainFunc(),
			NRet:    1,
			Protect: true, // If ``Protect`` is false,
			// GopherLua will panic instead of returning an ``error`` value.
			Handler: &lua.LFunction{
				GFunction: func(*lua.LState) int {
					return 0
				},
			},
		}, lua.LBool(false))
		if err != nil {
			glogger.GLogger.Error("app.VM().CallByParam error:", err)
			return
		}
	}(ctx)
	glogger.GLogger.Info("App started:", app.UUID)
	return nil
}

/*
*
* 从内存里面删除APP
*
 */
func (as *AppStack) RemoveApp(uuid string) error {
	if app, ok := as.applications[uuid]; ok {
		app.Stop()
		delete(as.applications, uuid)
	}
	glogger.GLogger.Info("App removed:", uuid)
	return nil
}

/*
*
* 停止应用并不删除应用, 将其进程结束，状态置0
*
 */
func (as *AppStack) StopApp(uuid string) error {
	if app, ok := as.applications[uuid]; ok {
		app.Stop()
		app.AppState = 0
	}
	glogger.GLogger.Info("App stopped:", uuid)
	return nil
}

/*
*
* 更新应用信息
*
 */
func (as *AppStack) UpdateApp(app typex.Application) error {
	if oldApp, ok := as.applications[app.UUID]; ok {
		oldApp.Name = app.Name
		oldApp.Version = app.Version
		glogger.GLogger.Info("App updated:", app.UUID)
		return nil
	}
	return fmt.Errorf("update failed, app not exists:%s", app.UUID)

}
func (as *AppStack) GetApp(uuid string) *typex.Application {
	if app, ok := as.applications[uuid]; ok {
		return app
	}
	return nil
}

/*
*
* 获取列表
*
 */
func (as *AppStack) ListApp() []*typex.Application {
	apps := []*typex.Application{}
	for _, v := range as.applications {
		apps = append(apps, v)
	}
	return apps
}

func (as *AppStack) Stop() {
	for _, app := range as.applications {
		glogger.GLogger.Info("Stop App:", app.UUID)
		app.Stop()
		glogger.GLogger.Info("Stop App:", app.UUID, " Successfully")
	}
	glogger.GLogger.Info("Appstack stopped")

}
func (as *AppStack) GetRuleX() typex.RuleX {
	return as.re
}
