package engine

import (
	"github.com/i4de/rulex/typex"
)

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