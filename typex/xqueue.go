package typex

import (
	"context"
	"errors"
	"fmt"

	"github.com/hootrhino/rulex/glogger"
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
	Debug bool // 是否是Debug消息
	I     *InEnd
	O     *OutEnd
	D     *Device
	E     RuleX
	Data  string
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
						// 如果是Debug消息直接打印出来
						qd.E.RunSourceCallbacks(qd.I, qd.Data)
					}
					if qd.D != nil {
						qd.E.RunDeviceCallbacks(qd.D, qd.Data)
						qd.E.RunHooks(qd.Data)
					}
					if qd.O != nil {
						v, ok := qd.E.AllOutEnd().Load(qd.O.UUID)
						if ok {
							target := v.(*OutEnd).Target
							if target == nil {
								continue
							}
							if _, err := target.To(qd.Data); err != nil {
								glogger.GLogger.Error(err)
								qd.E.GetMetricStatistics().IncOutFailed()
							} else {
								qd.E.GetMetricStatistics().IncOut()
							}
						}
					}
				}
			}
		}
	}(GCTX, DefaultDataCacheQueue)
}
