package rulexlib

import (
	"errors"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) error {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		e.PushOutQueue(outEnd, incoming)
		return nil
	} else {
		msg := "target not found:" + uuid
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	}
}
