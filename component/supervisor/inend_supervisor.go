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

package supervisor

import (
	"context"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/component/internotify"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

func StartInSupervisor(ctx context.Context, in *typex.InEnd, ruleEngine typex.RuleX) {
	UUID := in.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				glogger.GLogger.Debugf("Source Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-typex.GCTX.Done():
			{
				return
			}
		default:
			{
			}
		}
		currentIn := ruleEngine.GetInEnd(UUID)
		if currentIn == nil {
			glogger.GLogger.Debugf("Source:%v Deleted, supervisor exit", UUID)
			return
		}
		if currentIn.Source.Status() == typex.SOURCE_STOP {
			glogger.GLogger.Debugf("Source:%v Stopped, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentIn.Source.Status() == typex.SOURCE_DOWN {
			info := fmt.Sprintf("Source:%v DOWN, supervisor try to Restart", UUID)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  "SOURCE",
				Event: "event.down",
				Ts:    uint64(time.Now().UnixNano()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			// go LoadNewestInEnd(UUID, ruleEngine)
			return
		}
		// glogger.GLogger.Debugf("Supervisor Get Source :%v state:%v", UUID, currentIn.Source.Status().String())
		<-ticker.C
	}
}
