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
	"encoding/json"

	"github.com/hootrhino/rulex/typex"
)

// Default InteractQueue
var __DefaultInteractQueue InteractQueue

/*
*
* 前后端交互数据
*
 */
type InteractQueueData struct {
	Topic       string                 `json:"topic"`
	ComponentId string                 `json:"componentId"`
	Data        map[string]interface{} `json:"data"`
}

func (v InteractQueueData) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

type InteractQueue struct {
	inQueue  chan InteractQueueData // 外面的数据进来后进此管道
	outQueue chan InteractQueueData // 任何发到这个管道的数据都会被发到外部Pipe
	rulex    typex.RuleX
}

func InitInteractQueue(rulex typex.RuleX, maxQueueSize int) *InteractQueue {
	__DefaultInteractQueue = InteractQueue{
		inQueue:  make(chan InteractQueueData, maxQueueSize),
		outQueue: make(chan InteractQueueData, maxQueueSize),
		rulex:    rulex,
	}
	return &__DefaultInteractQueue
}

/*
*
* 给前端推送数据
*
 */
func SendData(data InteractQueueData) {
	__DefaultInteractQueue.outQueue <- data
}
func ReceiveData(data InteractQueueData) {
	__DefaultInteractQueue.inQueue <- data
}

/*
*
* GetQueue
*
 */
func InQueue() chan InteractQueueData {
	return __DefaultInteractQueue.inQueue
}
func OutQueue() chan InteractQueueData {
	return __DefaultInteractQueue.outQueue
}
