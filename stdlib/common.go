package stdlib

import (
	"rulex/statistics"
	"rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) {

	statistics.IncOut()
	outEnds := e.AllOutEnd()
	outEnd, _ := outEnds.Load(uuid)
	e.PushQueue(typex.QueueData{
		In:   nil,
		Out:  outEnd.(*typex.OutEnd),
		E:    e,
		Data: incoming,
	})

}
