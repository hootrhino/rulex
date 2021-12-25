package typex

import (
	"github.com/cjoudrey/gluaurl"
	luajson "github.com/wwhai/gopher-json"

	lua "github.com/yuin/gopher-lua"
)

type RuleStatus int

//
//  规则状态：
//  0: 停止
//  1: 运行中
const RULE_STOP RuleStatus = 0
const RULE_RUNNING RuleStatus = 1

//
// 规则描述
//
type Rule struct {
	Id          string      `json:"id"`
	UUID        string      `json:"uuid"`
	Status      RuleStatus  `json:"status"`
	Name        string      `json:"name"`
	From        []string    `json:"from"`
	Actions     string      `json:"actions"`
	Success     string      `json:"success"`
	Failed      string      `json:"failed"`
	Description string      `json:"description"`
	VM          *lua.LState `json:"-"`
}

//
// New
//
func NewRule(e RuleX,
	uuid string,
	name string,
	description string,
	from []string,
	success string,
	actions string,
	failed string) *Rule {
	return &Rule{
		UUID:        uuid,
		Name:        name,
		Description: description,
		From:        from,
		Status:      RULE_RUNNING, // 默认为启用
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		VM: lua.NewState(lua.Options{
			RegistrySize:     1024 * 1024,
			RegistryMaxSize:  1024 * 1024,
			RegistryGrowStep: 32,
		}),
	}
}

/*
*
* 配置LUA虚拟机
*
 */
func (r *Rule) SetVM(o lua.Options) {
	r.VM.Options = o
}
func (r *Rule) LoadLib(rx RuleX, lib XLib) {
	// log.Info("LoadLib:", lib.Name())
	rulex := r.VM.G.Global
	//
	// rulexlib: 标准库命名空间
	//
	r.VM.SetGlobal("rulexlib", rulex)
	r.VM.PreloadModule("json", luajson.Loader)
	r.VM.PreloadModule("url", gluaurl.Loader)
	//
	mod := r.VM.SetFuncs(rulex, map[string]lua.LGFunction{
		lib.Name(): lib.LibFun(rx),
	})
	r.VM.Push(mod)
}
