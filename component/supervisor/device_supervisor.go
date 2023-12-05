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

func StartDeviceSupervisor(ctx context.Context, device *typex.Device, ruleEngine typex.RuleX) {
	UUID := device.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				glogger.GLogger.Debugf("Device Context cancel:%v, supervisor exit", UUID)
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
		currentDevice := ruleEngine.GetDevice(UUID)
		if currentDevice == nil {
			glogger.GLogger.Debugf("Device:%v Deleted, supervisor exit", UUID)
			return
		}
		if currentDevice.Device.Status() == typex.DEV_STOP {
			glogger.GLogger.Debugf("Device:%v Stopped, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentDevice.Device.Status() == typex.DEV_DOWN {
			info := fmt.Sprintf("Device:%v DOWN, supervisor try to Restart", UUID)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  "DEVICE",
				Event: "event.down",
				Ts:    uint64(time.Now().UnixNano()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			// go LoadNewestDevice(UUID, ruleEngine)
			return
		}
		// glogger.GLogger.Debugf("Supervisor Get Device :%v state:%v", UUID, currentDevice.Device.Status().String())
		<-ticker.C
	}
}
