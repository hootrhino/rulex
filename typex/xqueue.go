package typex

import (
	"context"
	"errors"
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/statistics"
)

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
	// glogger.GLogger.Debug("DataCacheQueue Push:", d.Data)
	// 比较数据和容积
	if len(q.Queue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v", q.GetSize(), len(q.Queue)+1)
		glogger.GLogger.Error(msg)
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

// 此处内置的消息队列用了go的channel, 看似好像很简单，但是经过测试发现完全满足网关需求，甚至都性能过剩了
// 因此大家看到这里务必担心, 我也知道有很精美的高级框架, 但是用简单的方法来实现功能不是更好吗？
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
						// glogger.GLogger.Debug("RunDeviceCallbacks Device:", qd.D.UUID)
						qd.E.RunDeviceCallbacks(qd.D, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.O != nil {
						v, ok := qd.E.AllOutEnd().Load(qd.O.UUID)
						if ok {
							if v.(*OutEnd).Target == nil {
								continue
							}
							if _, err := v.(*OutEnd).Target.To(qd.Data); err != nil {
								glogger.GLogger.Error(err)
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
