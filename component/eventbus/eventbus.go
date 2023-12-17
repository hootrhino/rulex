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

package eventbus

import (
	"context"
	"fmt"

	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

var __DefaultEventBus *EventBus

/*
*
* 内部消息总线
*
 */
type EventMessage struct {
	Payload string
}

func (E EventMessage) String() string {
	return fmt.Sprintf("Event Message@ Payload: %s", E.Payload)
}

type Topic struct {
	Topic       string
	channel     chan EventMessage
	Subscribers map[string]*Subscriber
	ctx         context.Context
	cancel      context.CancelFunc
}
type Subscriber struct {
	id       string
	Callback func(Topic string, Msg EventMessage)
}
type EventBus struct {
	// Topic, chan EventMessage
	// 给每个订阅者分配一个Channel，实现消息订阅
	// Topic一样的会挂在同一个树上
	Topics map[string]*Topic // 订阅树: MAP<Topic>[]Subscribers
}

func InitEventBus() *EventBus {
	__DefaultEventBus = &EventBus{
		Topics: map[string]*Topic{},
	}
	return __DefaultEventBus
}

/*
*
* 订阅
*
 */
func Subscribe(topic string, subscribe *Subscriber) {
	NewUUID := utils.MakeUUID("SUB")
	subscribe.id = NewUUID
	var T *Topic
	Ok := false
	if T, Ok = __DefaultEventBus.Topics[topic]; Ok {
		T.Subscribers[subscribe.id] = subscribe
	} else {
		T = new(Topic)
		T.channel = make(chan EventMessage, 100)
		T.Subscribers = map[string]*Subscriber{}
		T.Topic = topic
		T.Subscribers[subscribe.id] = subscribe
		__DefaultEventBus.Topics[topic] = T
		ctx, cancel := context.WithCancel(typex.GCTX)
		T.ctx = ctx
		T.cancel = cancel
		go func(T *Topic) {
			for {
				select {
				case <-T.ctx.Done():
					{
						return
					}
				case Msg := <-T.channel:
					for _, Subscriber := range T.Subscribers {
						if Subscriber.Callback != nil {
							Subscriber.Callback(T.Topic, Msg)
						}
					}
				}
			}

		}(T)
	}

}

/*
*
* 取消订阅
*
 */
func UnSubscribe(topic string, subscribe Subscriber) {
	T, Ok1 := __DefaultEventBus.Topics[topic]
	if Ok1 {
		if _, Ok2 := T.Subscribers[subscribe.id]; Ok2 {
			delete(__DefaultEventBus.Topics[topic].Subscribers, subscribe.id)
		}
	}
	// 当没有订阅者的时候直接删除这个Topic
	if len(T.Subscribers) == 0 {
		T.cancel()
		delete(__DefaultEventBus.Topics, topic)
	}
}

/*
*
* 发布
*
 */
func Publish(topic string, Msg EventMessage) {
	T, Ok1 := __DefaultEventBus.Topics[topic]
	if Ok1 {
		T.channel <- Msg
	}
}

/*
*
* 释放所有
*
 */
func Flush() {
	for _, T := range __DefaultEventBus.Topics {
		for _, S := range T.Subscribers {
			delete(T.Subscribers, S.id)
		}
		T.cancel()
		delete(__DefaultEventBus.Topics, T.Topic)
	}
}
