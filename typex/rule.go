package typex

import (
	lua "github.com/hootrhino/gopher-lua"
)

type RuleStatus int

const _VM_Registry_Size int = 1024 * 1024    // 默认堆栈大小
const _VM_Registry_MaxSize int = 1024 * 1024 // 默认最大堆栈
const _VM_Registry_GrowStep int = 32         // 默认CPU消耗

// 规则状态：
// 0: 停止
// 1: 运行中
const RULE_STOP RuleStatus = 0
const RULE_RUNNING RuleStatus = 1

// 规则描述
type Rule struct {
	Id          string      `json:"id"`
	UUID        string      `json:"uuid"`
	Type        string      `json:"type"` // 脚本类型，目前支持"lua"
	Status      RuleStatus  `json:"status"`
	Name        string      `json:"name"`
	FromSource  []string    `json:"fromSource"` // 来自数据源
	FromDevice  []string    `json:"fromDevice"` // 来自设备
	Actions     string      `json:"actions"`
	Success     string      `json:"success"`
	Failed      string      `json:"failed"`
	Description string      `json:"description"`
	LuaVM       *lua.LState `json:"-"` // Lua VM
}

func NewLuaRule(e RuleX,
	uuid string,
	name string,
	description string,
	fromSource []string,
	fromDevice []string,
	success string,
	actions string,
	failed string) *Rule {
	rule := NewRule(e,
		uuid,
		name,
		description,
		fromSource,
		fromDevice,
		success,
		actions,
		failed)
	return rule
}

// New
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
		Type:        "lua", // 默认执行lua脚本
		Description: description,
		FromSource:  fromSource,
		FromDevice:  fromDevice,
		Status:      RULE_RUNNING, // 默认为启用
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		LuaVM: lua.NewState(lua.Options{
			RegistrySize:     _VM_Registry_Size,
			RegistryMaxSize:  _VM_Registry_MaxSize,
			RegistryGrowStep: _VM_Registry_GrowStep,
		}),
	}
}

/*
*
* 加载外部LUA脚本，方便用户自己写一些东西
* 需要注意的：
* - 不要和标准库里面的变量冲突了
* - 默认加载到 _G 环境里
 */
func (r *Rule) LoadExternLuaLib(path string) error {
	return r.LuaVM.DoFile(path)
}

/*
*
* AddLib: 根据 KV形式加载库(推荐)
*  - Global: 命名空间
*   - funcName: 函数名称
 */
func (r *Rule) AddLib(rx RuleX, Global string, funcName string,
	f func(l *lua.LState) int) {
	rulexTb := r.LuaVM.G.Global
	r.LuaVM.SetGlobal(Global, rulexTb)
	loadLib(rulexTb, r.LuaVM, funcName, f)
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
