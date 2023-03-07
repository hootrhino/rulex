package typex

import (
	"context"

	lua "github.com/yuin/gopher-lua"
)

// lua 虚拟机的参数
const p_VM_Registry_Size int = 1024 * 1024    // 默认堆栈大小
const p_VM_Registry_MaxSize int = 1024 * 1024 // 默认最大堆栈
const p_VM_Registry_GrowStep int = 32         // 默认CPU消耗
type AppState int

/*
*
* 轻量级应用
*
 */
type Application struct {
	UUID        string             `json:"uuid"`     // 名称
	Name        string             `json:"name"`     // 名称
	Version     string             `json:"version"`  // 版本号
	AppState    AppState           `json:"appState"` // 状态: 1 运行中, 0 停止
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
		RegistrySize:     p_VM_Registry_Size,
		RegistryMaxSize:  p_VM_Registry_MaxSize,
		RegistryGrowStep: p_VM_Registry_GrowStep,
	})
	return app
}

func (app *Application) SetCnC(ctx context.Context, cancel context.CancelFunc) {
	app.ctx = ctx
	app.cancel = cancel
}
func (app *Application) GetCnC() (context.Context, context.CancelFunc) {
	return app.ctx, app.cancel
}

func (app *Application) SetMainFunc(f *lua.LFunction) {
	app.luaMainFunc = f
}
func (app *Application) GetMainFunc() *lua.LFunction {
	return app.luaMainFunc
}

func (app *Application) VM() *lua.LState {
	return app.vm
}
func (app *Application) Release() {
	app.cancel()
	app.vm.Close()
	app.vm = nil
}

/*
*
* APP Stack 管理器
*
 */
type XAppStack interface {
	GetRuleX() RuleX
	ListApp() []*Application
	LoadApp(app *Application) error
	GetApp(uuid string) *Application
	RemoveApp(uuid string) error
	UpdateApp(app *Application) error
	StopApp(uuid string) error
	Stop()
}
