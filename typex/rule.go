package typex

import (
	lua "github.com/yuin/gopher-lua"
)

type RuleStatus int

const _VM_Registry_Size int = 1024 * 1024    // 默认堆栈大小
const _VM_Registry_MaxSize int = 1024 * 1024 // 默认最大堆栈
const _VM_Registry_GrowStep int = 32         // 默认CPU消耗

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
	FromSource  []string    `json:"fromSource"` // 来自数据源
	FromDevice  []string    `json:"fromDevice"` // 来自设备
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
	fromSource []string,
	fromDevice []string,
	success string,
	actions string,
	failed string) *Rule {
	return &Rule{
		UUID:        uuid,
		Name:        name,
		Description: description,
		FromSource:  fromSource,
		FromDevice:  fromDevice,
		Status:      RULE_RUNNING, // 默认为启用
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		VM: lua.NewState(lua.Options{
			RegistrySize:     _VM_Registry_Size,
			RegistryMaxSize:  _VM_Registry_MaxSize,
			RegistryGrowStep: _VM_Registry_GrowStep,
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

/*
*
* AddLib: 根据 KV形式加载库(推荐)
*  - Global: 命名空间
*   - funcName: 函数名称
 */
func (r *Rule) AddLib(rx RuleX, Global string, funcName string, f func(*lua.LState) int) {
	rulexTb := r.VM.G.Global
	r.VM.SetGlobal(Global, rulexTb)
	loadLib(rulexTb, r.VM, funcName, f)
}

func loadLib(
	tb *lua.LTable,
	VM *lua.LState,
	funcName string,
	f func(*lua.LState) int,
) {
	mod := VM.SetFuncs(tb, map[string]lua.LGFunction{
		funcName: f,
	})
	VM.Push(mod)
}
