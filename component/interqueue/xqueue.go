package interqueue

import (
	"context"
	"errors"
	"fmt"

	"github.com/hootrhino/rulex/component/intermetric"
	"github.com/hootrhino/rulex/glogger"
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
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.Queue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.Queue)+1)
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

// TODO: 下个版本更换为可扩容的Chan
func StartDataCacheQueue() {

	go func(ctx context.Context, xQueue XQueue) {
		for {
			select {
			case <-ctx.Done():
				return
			// 这个地方不能阻塞，需要借助一个外部queue
			// push qd -> Queue
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
		glogger.GLogger.Error("PushInQueue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
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
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
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
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
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
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
}
