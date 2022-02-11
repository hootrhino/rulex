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
	In   *InEnd
	Out  *OutEnd
	E    RuleX
	Data string
}

func (qd QueueData) String() string {
	return "QueueData@In:" + qd.In.UUID + ", Data:" + qd.Data
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
					// 消息队列有2种用法:
					// 1 进来的数据缓存
					// 2 出去的消息缓存
					// 只需要判断 in 或者 out 是不是 nil即可
					//
					if qd.In != nil {
						qd.E.RunLuaCallbacks(qd.In, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.Out != nil {
						outEnds := qd.E.AllOutEnd()
						v, ok := outEnds.Load(qd.Out.UUID)
						if ok {
							if err := v.(*OutEnd).Target.To(qd.Data); err != nil {
								statistics.IncOut()
							} else {
								statistics.IncOutFailed()
							}

						}
					}
				}
			}
		}
	}(GCTX, DefaultDataCacheQueue)
}
