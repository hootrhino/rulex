package engine

import (
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
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
	if err := e.AppStack.StopApp(uuid); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}
func (e *RuleEngine) RemoveApp(uuid string) error {
	if err := e.AppStack.RemoveApp(uuid); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}
func (e *RuleEngine) LoadApp(app *typex.Application) error {
	if err := e.AppStack.LoadApp(app); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (e *RuleEngine) StartApp(uuid string) error {
	if err := e.AppStack.StartApp(uuid); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}
