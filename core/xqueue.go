/*
*
* Data Cache Queue
*
 */
package core

import (
	"context"
	"errors"
	"fmt"
	"rulex/typex"

	"github.com/ngaut/log"
)

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
	In   *typex.InEnd
	E    typex.RuleX
	Data string
}

func (qd QueueData) String() string {
	return "QueueData@In:" + qd.In.Id + ", Data:" + qd.Data
}

var DefaultDataCacheQueue XQueue

/*
*
* NewXQueue
*
 */

func InitXQueue(size int, rulex typex.RuleX) {
	log.Info("Init XQueue max size is:", size)
	DefaultDataCacheQueue = &DataCacheQueue{
		Size:  size,
		Queue: make(chan QueueData, size),
	}
	go func(ctx context.Context, xQueue XQueue) {
		for {
			log.Info("Size is: ", xQueue.GetSize())
			select {

			case qd := <-xQueue.GetQueue():
				qd.E.RunLuaCallbacks(qd.In, qd.Data)
				qd.E.RunHooks(qd.Data)
			}
		}

	}(context.Background(), DefaultDataCacheQueue)
}

/*
*
* DataCacheQueue
*
 */
type DataCacheQueue struct {
	Size  int
	Queue chan QueueData
}

func (q *DataCacheQueue) GetSize() int {
	return q.Size
}

/*
*
* Push
*
 */
func (q *DataCacheQueue) Push(d QueueData) error {
	if len(q.Queue)+1 > q.Size {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v", q.Size, len(q.Queue)+1)
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
