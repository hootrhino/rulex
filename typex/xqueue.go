package typex

import (
	"errors"
	"fmt"

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
	return "QueueData@In:" + qd.In.Id + ", Data:" + qd.Data
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
