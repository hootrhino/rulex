package rulexlib

import (
	"errors"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/component/interqueue"
	"github.com/hootrhino/rulex/typex"
)

func handleDataFormat(e typex.RuleX, uuid string, incoming string) error {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		return interqueue.DefaultDataCacheQueue.PushOutQueue(outEnd, incoming)
	}
	msg := "target not found:" + uuid
	glogger.GLogger.Error(msg)
	return errors.New(msg)

}
