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
	GetInQueue() chan QueueData
	GetOutQueue() chan QueueData
	GetDeviceQueue() chan QueueData
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
	Queue       chan QueueData
	OutQueue    chan QueueData
	InQueue     chan QueueData
	DeviceQueue chan QueueData
	rulex       typex.RuleX
}

func InitDataCacheQueue(rulex typex.RuleX, maxQueueSize int) XQueue {
	DefaultDataCacheQueue = &DataCacheQueue{
		Queue:       make(chan QueueData, maxQueueSize),
		OutQueue:    make(chan QueueData, maxQueueSize),
		InQueue:     make(chan QueueData, maxQueueSize),
		DeviceQueue: make(chan QueueData, maxQueueSize),
		rulex:       rulex,
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

/*
*
* GetQueue
*
 */
func (q *DataCacheQueue) GetInQueue() chan QueueData {
	return q.InQueue
}

/*
*
* GetQueue
*
 */
func (q *DataCacheQueue) GetOutQueue() chan QueueData {
	return q.OutQueue
}

/*
*
*GetDeviceQueue
*
 */
func (q *DataCacheQueue) GetDeviceQueue() chan QueueData {
	return q.DeviceQueue
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
			case qd := <-xQueue.GetInQueue():
				{
					if qd.I != nil {
						qd.E.RunSourceCallbacks(qd.I, qd.Data)
					}
				}
			case qd := <-xQueue.GetDeviceQueue():
				{
					if qd.D != nil {
						qd.E.RunDeviceCallbacks(qd.D, qd.Data)
					}
				}
			case qd := <-xQueue.GetOutQueue():
				{
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
			case qd := <-xQueue.GetQueue(): // 马上废弃
				{
					if qd.I != nil {
						qd.E.RunSourceCallbacks(qd.I, qd.Data)
					}
					if qd.D != nil {
						qd.E.RunDeviceCallbacks(qd.D, qd.Data)
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

/*
*
*
*
 */
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

/*
*
*PushInQueue
*
 */
func (q *DataCacheQueue) PushInQueue(in *typex.InEnd, data string) error {
	qd := QueueData{
		E:    q.rulex,
		I:    in,
		O:    nil,
		Data: data,
	}
	err := q.pushIn(qd)
	if err != nil {
		glogger.GLogger.Error("Push InQueue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
}

/*
*
* PushDeviceQueue
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
	err := q.pushDevice(qd)
	if err != nil {
		glogger.GLogger.Error("Push Device Queue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
}

/*
*
* PushOutQueue
*
 */
func (q *DataCacheQueue) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := QueueData{
		E:    q.rulex,
		D:    nil,
		I:    nil,
		O:    out,
		Data: data,
	}
	err := q.pushOut(qd)
	if err != nil {
		glogger.GLogger.Error("Push OutQueue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
}

/*
*
* Push
*
 */
func (q *DataCacheQueue) pushIn(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.InQueue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.Queue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.InQueue <- d
		return nil
	}
}

/*
*
* Push
*
 */
func (q *DataCacheQueue) pushOut(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.OutQueue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.Queue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.OutQueue <- d
		return nil
	}
}

/*
*
* Push
*
 */
func (q *DataCacheQueue) pushDevice(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.DeviceQueue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.Queue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.DeviceQueue <- d
		return nil
	}
}
