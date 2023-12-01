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

package intereventbus

import (
	"context"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultInternalEventSource *InternalEventBus

// ---------------------------------------------------------
// Type
// ---------------------------------------------------------
type BaseEvent struct {
	Type  string
	Event string
	Ts    uint64
	Info  interface{}
}

type InternalEventBus struct {
	Queue chan BaseEvent
	Rulex typex.RuleX
}

func (q *InternalEventBus) GetSize() int {
	return cap(q.Queue)
}

/*
*
  - 内部事件，例如资源挂了或者设备离线、超时等等,该资源是单例模式,
    维护一个channel来接收各种事件，将收到的消息吐到InterQueue即可

*
*/
func InitInternalEventSource(r typex.RuleX) *InternalEventBus {
	__DefaultInternalEventSource = new(InternalEventBus)
	__DefaultInternalEventSource.Rulex = r
	return __DefaultInternalEventSource
}

/*
*
* 监控chan
*
 */
func StartInternalEventQueue() {
	go func(ctx context.Context, InternalEventBus InternalEventBus) {
		for {
			select {
			case <-ctx.Done():
				return
			case Event := <-__DefaultInternalEventSource.Queue:
				{
					glogger.GLogger.Debug("Event:", Event)

				}
			}
		}
	}(typex.GCTX, *__DefaultInternalEventSource)
}
