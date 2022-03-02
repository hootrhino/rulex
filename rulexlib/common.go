package rulexlib

import (
	"rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) {
	outEnd := e.GetOutEnd(uuid)
	e.PushOutQueue(outEnd, incoming)
}
