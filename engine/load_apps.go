// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
