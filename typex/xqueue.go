package typex

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
