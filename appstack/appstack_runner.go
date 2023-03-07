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
	Applications map[string]*typex.Application
}

func NewAppStack(re typex.RuleX) typex.XAppStack {
	as := new(AppStack)
	as.re = re
	as.Applications = map[string]*typex.Application{}
	return as
}

/*
*
* 加载本地文件到lua虚拟机
*
 */
func (as *AppStack) LoadApp(app *typex.Application) error {

	// 临时校验语法
	tempVm := lua.NewState()
	bytes, err := os.ReadFile("./apps/" + app.Filepath)
	if err != nil {
		return err
	}
	if err := tempVm.DoString(string(bytes)); err != nil {
		return err
	}
	// 检查名称
	AppNAME := tempVm.GetGlobal("AppNAME")
	if AppNAME == nil {
		return fmt.Errorf("'AppNAME' field not exists")
	}
	if AppNAME.Type() != lua.LTString {
		return fmt.Errorf("'AppNAME' must be string")
	}
	// 检查类型
	AppVERSION := tempVm.GetGlobal("AppVERSION")
	if AppVERSION == nil {
		return fmt.Errorf("'AppVERSION' field not exists")
	}
	if AppVERSION.Type() != lua.LTString {
		return fmt.Errorf("'AppVERSION' must be string")
	}
	// 检查描述信息
	AppDESCRIPTION := tempVm.GetGlobal("AppDESCRIPTION")
	if AppDESCRIPTION == nil {
		if AppDESCRIPTION.Type() != lua.LTString {
			return fmt.Errorf("'AppDESCRIPTION' must be string")
		}
	}

	// 检查函数入口
	AppMain := tempVm.GetGlobal("Main")
	if AppMain == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMain.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	// 释放语法验证阶段的临时虚拟机
	tempVm.Close()
	tempVm = nil
	//----------------------------------------------------------------------------------------------
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
	if err := as.startApp(app); err != nil {
		return err
	}
	// 加载到内存里
	as.Applications[app.UUID] = app

	return nil
}

/*
*
* 启动 function Main(args) --do-some-thing-- return 0 end
*
 */
func (as *AppStack) startApp(app *typex.Application) error {
	// args := lua.LBool(false) // Main的参数，未来准备扩展
	ctx, cancel := context.WithCancel(typex.GCTX)
	app.SetCnC(ctx, cancel)
	go func(ctx context.Context) {
		defer func() {
			glogger.GLogger.Debug("app exit:", app.UUID)
		}()
		app.VM().SetContext(ctx)
		err := app.VM().CallByParam(lua.P{
			Fn:      app.GetMainFunc(), // 回调函数
			NRet:    1,                 // 一个返回值
			Protect: true,              // 受保护
		}, lua.LBool(false))
		if err != nil {
			glogger.GLogger.Error("startApp error:", err)
			return
		}
	}(ctx)

	return nil
}

/*
*
* 删除APP
*
 */
func (as *AppStack) RemoveApp(uuid string) error {
	if app, ok := as.Applications[uuid]; ok {
		app.Release()
		delete(as.Applications, uuid)
	}
	return nil
}
func (as *AppStack) StopApp(uuid string) error {
	if app, ok := as.Applications[uuid]; ok {
		app.Release()
		app.AppState = 0
	}
	return nil
}

/*
*
* 更新应用信息
*
 */
func (as *AppStack) UpdateApp(app *typex.Application) error {
	if _, ok := as.Applications[app.UUID]; ok {
		as.Applications[app.UUID] = app
	}
	return fmt.Errorf("update failed, app not exists:%s", app.UUID)

}
func (as *AppStack) GetApp(uuid string) *typex.Application {
	if app, ok := as.Applications[uuid]; ok {
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
	for _, v := range as.Applications {
		apps = append(apps, v)
	}
	return apps
}

func (as *AppStack) Stop() {
	for _, app := range as.Applications {
		app.Release()
	}
}
func (as *AppStack) GetRuleX() typex.RuleX {
	return as.re
}
