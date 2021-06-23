package x

import (
	"rulenginex/extralib"

	luajson "github.com/wwhai/gopher-json"
	lua "github.com/yuin/gopher-lua"
)

//
//
type Resource interface {
	Start(e *RuleEngine, successCallBack func(), errorCallback func(error)) error
	Register(inEndId string) error
	Reload()
	Pause()
	Status() int
	Stop()
}

//
type InEnd struct {
	Id          string
	State       int
	Type        string
	Name        string
	Description string
	Binds       *map[string]Rule
	Config      *map[string]interface{}
}

//
//
//
func (in *InEnd) GetResource() {

}

//
//
//
type OutEnd struct {
	Id          string
	Type        string
	Name        string
	Description string
	Config      map[string]interface{}
}

// InEnd
//
//
type Rule struct {
	Id          string
	Name        string
	Description string
	VM          *lua.LState
	From        []string
	Actions     string
	Success     string
	Failed      string
}

//
// New
//
func NewRule(id string,
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
	extralib.LoadDecodeLib(vm)
	extralib.LoadEncodeLib(vm)
	extralib.LoadSqlLib(vm)
	luajson.Preload(vm)
	return &Rule{
		Id:          id,
		Name:        name,
		Description: description,
		From:        from,
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		VM:          vm,
	}
}

//
//
//
type StartCallback func()

//
//
//
type RuleEngine struct {
	ConfigMap *map[string]interface{}
}
