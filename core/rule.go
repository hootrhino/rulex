package core

import (
	luajson "github.com/wwhai/gopher-json"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
type rule struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	VM          *lua.LState `json:"-"`
	From        []string    `json:"from"`
	Actions     string      `json:"actions"`
	Success     string      `json:"success"`
	Failed      string      `json:"failed"`
}

//
// New
//
func NewRule(e *RuleEngine,
	name string,
	description string,
	from []string,
	success string,
	actions string,
	failed string) *rule {
	vm := lua.NewState(lua.Options{
		RegistrySize:     1024 * 20,
		RegistryMaxSize:  1024 * 80,
		RegistryGrowStep: 32,
	})
	LoadDbLib(e, vm)
	LoadJqLib(e, vm)
	luajson.Preload(vm)
	return &rule{
		Id:          MakeUUID("RULE"),
		Name:        name,
		Description: description,
		From:        from,
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		VM:          vm,
	}
}
