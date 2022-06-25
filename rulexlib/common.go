package rulexlib

import (
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		e.PushOutQueue(outEnd, incoming)
	}else {
		glogger.GLogger.Error()
	}
}
