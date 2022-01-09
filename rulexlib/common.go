package rulexlib

import (
	"rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) {
	outEnds := e.AllOutEnd()
	outEnd, _ := outEnds.Load(uuid)
	e.PushOutQueue(outEnd.(*typex.OutEnd), incoming)
}
