package appstack

import (
	"context"
	"fmt"
	"os"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	lua "github.com/yuin/gopher-lua"
)

// lua 虚拟机的参数
const _VM_Registry_Size int = 1024 * 1024    // 默认堆栈大小
const _VM_Registry_MaxSize int = 1024 * 1024 // 默认最大堆栈
const _VM_Registry_GrowStep int = 32         // 默认CPU消耗

/*
*
* 轻量级应用
*
 */
type Application struct {
	UUID        string             `json:"uuid"`     // 名称
	Name        string             `json:"name"`     // 名称
	Version     string             `json:"version"`  // 版本号
	Filepath    string             `json:"filepath"` // 文件路径, 是相对于main的apps目录
	luaMainFunc *lua.LFunction     `json:"-"`
	vm          *lua.LState        `json:"-"` // lua 环境
	ctx         context.Context    `json:"-"`
	cancel      context.CancelFunc `json:"-"`
}

func NewApplication(uuid, Name, Version, Filepath string) *Application {
	app := new(Application)
	app.Name = Name
	app.Version = Version
	app.Filepath = Filepath
	app.vm = lua.NewState(lua.Options{
		RegistrySize:     _VM_Registry_Size,
		RegistryMaxSize:  _VM_Registry_MaxSize,
		RegistryGrowStep: _VM_Registry_GrowStep,
	})
	return app
}

/*
*
* 管理器
*
 */
type AppStack struct {
	RuleEngine   typex.RuleX
	Applications map[string]*Application
}

func NewAppStack(rulex typex.RuleX) *AppStack {
	as := new(AppStack)
	as.Applications = map[string]*Application{}
	as.RuleEngine = rulex
	return as
}

/*
*
* 加载本地文件到lua虚拟机
*
 */
func (as *AppStack) LoadApp(app *Application) error {

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
	app.vm.DoString(string(bytes))
	// 检查函数入口
	AppMainVM := app.vm.GetGlobal("Main")
	if AppMainVM == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMainVM.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	// 抽取main
	fMain := *AppMainVM.(*lua.LFunction)
	app.luaMainFunc = &fMain
	// 加载库
	app.loadAppLib(as.RuleEngine)
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
func (as *AppStack) startApp(app *Application) error {
	// args := lua.LBool(false) // Main的参数，未来准备扩展
	ctx, cancel := context.WithCancel(typex.GCTX)
	app.ctx = ctx
	app.cancel = cancel
	go func(ctx context.Context) {
		defer func() {
			glogger.GLogger.Debug("app exit:", app.UUID)
		}()
		app.vm.SetContext(ctx)
		err := app.vm.CallByParam(lua.P{
			Fn:      app.luaMainFunc, // 回调函数
			NRet:    1,               // 一个返回值
			Protect: true,            // 受保护
		}, lua.LBool(false))
		if err != nil {
			glogger.GLogger.Error("startApp error:", err)
			return
		}
	}(app.ctx)

	return nil
}

/*
*
* 删除APP
*
 */
func (as *AppStack) RemoveApp(uuid string) error {
	if app, ok := as.Applications[uuid]; ok {
		app.cancel()
		app.vm.Close()
		delete(as.Applications, uuid)
		app.vm = nil
	}
	return nil
}

/*
*
* 更新应用信息
*
 */
func (as *AppStack) UpdateApp(app *Application) error {
	if _, ok := as.Applications[app.UUID]; ok {
		as.Applications[app.UUID] = app
	}
	return fmt.Errorf("update failed, app not exists:%s", app.UUID)

}

/*
*
* 获取列表
*
 */
func (as *AppStack) ListApp() []Application {
	apps := []Application{}
	for _, v := range as.Applications {
		apps = append(apps, *v)
	}
	return apps
}

func (as *AppStack) Stop() {
	for _, app := range as.Applications {
		app.cancel()
		app.vm.Close()
		app.vm = nil
	}
}
