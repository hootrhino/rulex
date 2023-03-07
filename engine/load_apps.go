package engine

import (
	"github.com/i4de/rulex/typex"
)

/*
*
* 都是一些对应用的CURD
*
 */
func (e *RuleEngine) AllApp() []*typex.Application {
	return e.AppStack.ListApp()
}
func (e *RuleEngine) GetApp(uuid string) *typex.Application {
	return e.AppStack.GetApp(uuid)
}
func (e *RuleEngine) StopApp(uuid string) error {
	return e.AppStack.StopApp(uuid)
}
func (e *RuleEngine) RemoveApp(uuid string) error {
	return e.AppStack.RemoveApp(uuid)
}
func (e *RuleEngine) LoadApp(app *typex.Application) error {
	return e.AppStack.LoadApp(app)
}

func (e *RuleEngine) StartApp(uuid string) error {
	return e.AppStack.StartApp(uuid)
}
