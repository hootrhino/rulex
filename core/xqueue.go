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

var DefaultDataCacheQueue typex.XQueue

/*
*
* NewXQueue
*
 */

func InitXQueue(size int, rulex typex.RuleX) {
	log.Info("Init XQueue max size is:", size)
	DefaultDataCacheQueue = &DataCacheQueue{
		Size:  size,
		Queue: make(chan typex.QueueData, size),
	}
	go func(ctx context.Context, xQueue typex.XQueue) {
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
	Queue chan typex.QueueData
}

func (q *DataCacheQueue) GetSize() int {
	return q.Size
}

/*
*
* Push
*
 */
func (q *DataCacheQueue) Push(d typex.QueueData) error {
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
func (q *DataCacheQueue) GetQueue() chan typex.QueueData {
	return q.Queue
}
