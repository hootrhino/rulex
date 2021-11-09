package typex

import (
	"errors"
	"reflect"

	"github.com/cjoudrey/gluaurl"
	luajson "github.com/wwhai/gopher-json"

	"rulex/utils"

	lua "github.com/yuin/gopher-lua"
)

//
//
//
type Rule struct {
	Id          string      `json:"id"`
	UUID        string      `json:"uuid"`
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
func NewRule(e RuleX,
	name string,
	description string,
	from []string,
	success string,
	actions string,
	failed string) *Rule {
	vm := lua.NewState(lua.Options{
		RegistrySize:     1024 * 20,
		RegistryMaxSize:  1024 * 80,
		RegistryGrowStep: 32,
	})

	return &Rule{
		UUID:        utils.MakeUUID("RULE"),
		Name:        name,
		Description: description,
		From:        from,
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		VM:          vm,
	}
}

// LUA Callback : Success
func (r *Rule) ExecuteSuccess() (interface{}, error) {
	return Execute(r.VM, "Success")
}

// LUA Callback : Failed

func (r *Rule) ExecuteFailed(arg lua.LValue) (interface{}, error) {
	return Execute(r.VM, "Failed", arg)
}

//
func (r *Rule) ExecuteActions(arg lua.LValue) (lua.LValue, error) {
	table := r.VM.GetGlobal("Actions")
	if table != nil && table.Type() == lua.LTTable {
		funcs := make(map[string]*lua.LFunction)
		table.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			t := reflect.TypeOf(f).Elem().Name()
			if t == "LFunction" {
				funcs[idx.String()] = f.(*lua.LFunction)
			}
		})
		return RunPipline(r.VM, funcs, arg)
	} else {
		return nil, errors.New("'Actions' not a lua table or not exist")
	}
}

func (r *Rule) LoadLib(rx RuleX, lib XLib) {
	// log.Info("LoadLib:", lib.Name())
	stdlib := r.VM.G.Global
	//
	r.VM.SetGlobal("stdlib", stdlib)
	r.VM.PreloadModule("json", luajson.Loader)
	r.VM.PreloadModule("url", gluaurl.Loader)
	//
	mod := r.VM.SetFuncs(stdlib, map[string]lua.LGFunction{
		lib.Name(): lib.LibFun(rx),
	})
	r.VM.Push(mod)
}
