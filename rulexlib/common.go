package rulexlib

import (
	"errors"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) error {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		return e.PushOutQueue(outEnd, incoming)
	}
	msg := "target not found:" + uuid
	glogger.GLogger.Error(msg)
	return errors.New(msg)

}
