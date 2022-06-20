package typex

import (
	"context"
	"errors"
	"fmt"
	"rulex/statistics"

	"github.com/ngaut/log"
)

//
//
//
var DefaultDataCacheQueue XQueue

/*
*
* XQueue
*
 */
type XQueue interface {
	GetQueue() chan QueueData
	GetSize() int
	Push(QueueData) error
}

//
type QueueData struct {
	I    *InEnd
	O    *OutEnd
	D    *Device
	E    RuleX
	Data string
}

func (qd QueueData) String() string {
	return "QueueData@In:" + qd.I.UUID + ", Data:" + qd.Data
}

/*
*
* NewXQueue
*
 */

/*
*
* DataCacheQueue
*
 */
type DataCacheQueue struct {
	Queue chan QueueData
}

func (q *DataCacheQueue) GetSize() int {
	return cap(q.Queue)
}

/*
*
* Push
*
 */
func (q *DataCacheQueue) Push(d QueueData) error {
	// 比较数据和容积
	if len(q.Queue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v", q.GetSize(), len(q.Queue)+1)
		log.Error(msg)
		return errors.New(msg)
	} else {
		q.Queue <- d
		return nil
	}
}

/*
*
* GetQueue
*
 */
func (q *DataCacheQueue) GetQueue() chan QueueData {
	return q.Queue
}

//
//
//
func StartQueue(maxQueueSize int) {
	DefaultDataCacheQueue = &DataCacheQueue{
		Queue: make(chan QueueData, maxQueueSize),
	}
	go func(ctx context.Context, xQueue XQueue) {
		for {
			select {
			case <-ctx.Done():
				return
			case qd := <-xQueue.GetQueue():
				{
					//
					// Rulex内置消息队列用法:
					// 1 进来的数据缓存
					// 2 出去的消息缓存
					// 3 设备数据缓存
					// 只需要判断 in 或者 out 是不是 nil即可
					//
					if qd.I != nil {
						qd.E.RunSourceCallbacks(qd.I, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.D != nil {
						qd.E.RunDeviceCallbacks(qd.D, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.O != nil {
						v, ok := qd.E.AllOutEnd().Load(qd.O.UUID)
						if ok {
							if _, err := v.(*OutEnd).Target.To(qd.Data); err != nil {
								log.Error(err)
								statistics.IncOutFailed()
							} else {
								statistics.IncOut()

							}
						}
					}
				}
			}
		}
	}(GCTX, DefaultDataCacheQueue)
}
