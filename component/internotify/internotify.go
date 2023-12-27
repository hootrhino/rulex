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

package internotify

import (
	"context"
	"errors"
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultInternalEventBus *InternalEventBus

// ---------------------------------------------------------
// Type
// ---------------------------------------------------------
// - SOURCE: 南向事件
// - DEVICE: 设备事件
// - TARGET: 北向事件
// - SYSTEM: 系统内部事件
// - HARDWARE: 硬件事件

type BaseEvent struct {
	Type  string
	Event string
	Ts    uint64
	Info  interface{}
}

func (be BaseEvent) String() string {
	return fmt.Sprintf(
		`BaseEvent@Type:%s, Event:%s, Ts:%d, Info:%s`,
		be.Type, be.Event, be.Ts, be.Info)

}

/*
*
* Push
*
 */
func Push(e BaseEvent) error {
	if len(__DefaultInternalEventBus.Queue)+1 > __DefaultInternalEventBus.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			__DefaultInternalEventBus.GetSize(), len(__DefaultInternalEventBus.Queue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		__DefaultInternalEventBus.Queue <- e
		return nil
	}
}

/*
*
* 内部事件总线
*
 */
type InternalEventBus struct {
	Queue       chan BaseEvent
	Rulex       typex.RuleX
	SourceCount uint
}

func (q *InternalEventBus) GetSize() int {
	return cap(q.Queue)
}
func RemoveSource() {
	if __DefaultInternalEventBus.SourceCount > 0 {
		__DefaultInternalEventBus.SourceCount--
	}
}
func AddSource() {
	__DefaultInternalEventBus.SourceCount++
}
func GetQueue() chan BaseEvent {
	return __DefaultInternalEventBus.Queue
}

/*
*
  - 内部事件，例如资源挂了或者设备离线、超时等等,该资源是单例模式,
    维护一个channel来接收各种事件，将收到的消息吐到InterQueue即可

*
*/
func InitInternalEventBus(r typex.RuleX, MaxQueueSize int) *InternalEventBus {
	__DefaultInternalEventBus = new(InternalEventBus)
	__DefaultInternalEventBus.Queue = make(chan BaseEvent, 1024)
	__DefaultInternalEventBus.Rulex = r
	return __DefaultInternalEventBus
}

/*
*
* 监控chan
*
 */
func StartInternalEventQueue() {
	go func(ctx context.Context, InternalEventBus *InternalEventBus) {
		for {
			// 当无订阅者时，及时释放channel里面的数据
			if __DefaultInternalEventBus.SourceCount == 0 {
				select {
				case <-ctx.Done():
					return
				case Event := <-InternalEventBus.Queue:
					{
						//
						// TODO 内部事件应该写入数据库, 主要是起通知作用
						//
						glogger.GLogger.Debug("Internal Event:", Event)

					}
				}
			}
		}
	}(typex.GCTX, __DefaultInternalEventBus)
}
