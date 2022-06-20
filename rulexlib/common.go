package rulexlib

import (
	"rulex/typex"

	"github.com/ngaut/log"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		e.PushOutQueue(outEnd, incoming)
	}else {
		log.Error()
	}
}
