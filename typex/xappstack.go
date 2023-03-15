package typex

import (
	"context"
	"log"
	"runtime"

	lua "github.com/i4de/gopher-lua"
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
	UUID        string             `json:"uuid"`      // 名称
	Name        string             `json:"name"`      // 名称
	Version     string             `json:"version"`   // 版本号
	AutoStart   bool               `json:"autoStart"` // 自动启动
	AppState    AppState           `json:"appState"`  // 状态: 1 运行中, 0 停止
	Filepath    string             `json:"filepath"`  // 文件路径, 是相对于main的apps目录
	luaMainFunc *lua.LFunction     `json:"-"`
	vm          *lua.LState        `json:"-"` // lua 环境
	ctx         context.Context    `json:"-"`
	cancel      context.CancelFunc `json:"-"`
}

func NewApplication(uuid, Name, Version, Filepath string) *Application {
	app := new(Application)
	app.Name = Name
	app.UUID = uuid
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

/*
*
* 源码bug，没有等字节码执行结束就直接给释放stack了，问题处在state.go:1391, 已经给作者提了issue，
* 如果1个月内不解决，准备自己fork一个过来维护.
* Issue: https://github.com/i4de/gopher-lua/discussions/430
 */
func (app *Application) Stop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[gopher-lua] app Stop:", app.UUID, ", with recover error: ", err)
		}
	}()
	app.vm.DoString(`function __() end __()`)
	app.vm.SetTop(-1)
	app.vm.Dead = true
	app.cancel()
	app.vm.Close() // app.vm.Close 会导致panic, 暂时先不处理
}

/*
*
* 清理内存
*
 */
func (app *Application) Remove() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[gopher-lua] app Remove:", app.UUID, ", with recover error: ", err)
		}
	}()
	app.Stop()
	app.vm = nil
	runtime.GC()
}

/*
*
* APP Stack 管理器
*
 */
type XAppStack interface {
	GetRuleX() RuleX
	ListApp() []*Application
	// 把配置里的应用信息加载到内存里
	LoadApp(app *Application) error
	GetApp(uuid string) *Application
	RemoveApp(uuid string) error
	UpdateApp(app Application) error
	// 启动一个停止的进程
	StartApp(uuid string) error
	StopApp(uuid string) error
	Stop()
}
