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

package interqueue

import (
	"context"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

// Default InteractQueue
var DefaultInteractQueue InteractQueue

/*
*
* 前后端交互数据
*
 */
type InteractQueueData struct {
	Topic string      `json:"topic"`
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
}

/*
*
* 前后端交互总线
*
 */
type XInteract interface {
	PushData(InteractQueueData)
	InQueue() chan InteractQueueData
	OutQueue() chan InteractQueueData
}
type InteractQueue struct {
	inQueue  chan InteractQueueData // 外面的数据进来后进此管道
	outQueue chan InteractQueueData // 任何发到这个管道的数据都会被发到外部Pipe
	rulex    typex.RuleX
}

func InitInteractQueue(rulex typex.RuleX, maxQueueSize int) XInteract {
	DefaultInteractQueue = InteractQueue{
		inQueue:  make(chan InteractQueueData, maxQueueSize),
		outQueue: make(chan InteractQueueData, maxQueueSize),
		rulex:    rulex,
	}
	return &DefaultInteractQueue
}

/*
*
* 给前端推送数据
*
 */
func (iq InteractQueue) PushData(data InteractQueueData) {

}

/*
*
* GetQueue
*
 */
func (iq InteractQueue) InQueue() chan InteractQueueData {
	return iq.inQueue
}
func (iq InteractQueue) OutQueue() chan InteractQueueData {
	return iq.outQueue
}

/*
*
* 启动双管道
*
 */
func StartInteractQueue() {
	// 监听管道
	go func(ctx context.Context, Interact XInteract) {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-Interact.InQueue():
				{
					glogger.GLogger.Debug("InQueue:", d)
				}
			case d := <-Interact.OutQueue():
				{
					glogger.GLogger.Debug("OutQueue:", d)
				}
			}
		}
	}(typex.GCTX, DefaultInteractQueue)
}
