package interqueue

import (
	"context"
	"errors"
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/component/intermetric"
	"github.com/hootrhino/rulex/typex"
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
	PushQueue(QueueData) error
	PushInQueue(in *typex.InEnd, data string) error
	PushOutQueue(in *typex.OutEnd, data string) error
	PushDeviceQueue(in *typex.Device, data string) error
}

type QueueData struct {
	Debug bool // 是否是Debug消息
	I     *typex.InEnd
	O     *typex.OutEnd
	D     *typex.Device
	E     typex.RuleX
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
	rulex typex.RuleX
}

func InitDataCacheQueue(rulex typex.RuleX, maxQueueSize int) XQueue {
	DefaultDataCacheQueue = &DataCacheQueue{
		Queue: make(chan QueueData, maxQueueSize),
		rulex: rulex,
	}
	return DefaultDataCacheQueue
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
func StartDataCacheQueue() {

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
							target := v.(*typex.OutEnd).Target
							if target == nil {
								continue
							}
							if _, err := target.To(qd.Data); err != nil {
								glogger.GLogger.Error(err)
								intermetric.IncOutFailed()
							} else {
								intermetric.IncOut()
							}
						}
					}
				}
			}
		}
	}(typex.GCTX, DefaultDataCacheQueue)
}

func (q *DataCacheQueue) PushQueue(qd QueueData) error {
	err := DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushQueue error:", err)
		// q.rulex.MetricStatistics.IncInFailed()
	} else {
		// e.MetricStatistics.IncIn()
	}
	return err
}
func (q *DataCacheQueue) PushInQueue(in *typex.InEnd, data string) error {
	qd := QueueData{
		E:    q.rulex,
		I:    in,
		O:    nil,
		Data: data,
	}
	err := DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushInQueue error:", err)
		// e.MetricStatistics.IncInFailed()
	} else {
		// e.MetricStatistics.IncIn()
	}
	return err
}

/*
*
* 设备数据入流引擎
*
 */
func (q *DataCacheQueue) PushDeviceQueue(Device *typex.Device, data string) error {
	qd := QueueData{
		D:    Device,
		E:    q.rulex,
		I:    nil,
		O:    nil,
		Data: data,
	}
	err := DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushInQueue error:", err)
		// q.rulex.MetricStatistics.IncInFailed()
	} else {
		// q.rulex.MetricStatistics.IncIn()
	}
	return err
}
func (q *DataCacheQueue) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := QueueData{
		E:    q.rulex,
		D:    nil,
		I:    nil,
		O:    out,
		Data: data,
	}
	err := DefaultDataCacheQueue.Push(qd)
	if err != nil {
		glogger.GLogger.Error("PushOutQueue error:", err)
		// e.MetricStatistics.IncInFailed()
	} else {
		// e.MetricStatistics.IncIn()
	}
	return err
}
