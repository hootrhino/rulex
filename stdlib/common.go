package stdlib

import (
	"encoding/json"
	"rulex/statistics"
	"rulex/typex"

	"github.com/ngaut/log"
)


func handleDataFormat(e typex.RuleX, id string, incoming string) {
	data := &map[string]interface{}{}
	err := json.Unmarshal([]byte(incoming), data)
	if err != nil {
		statistics.IncOutFailed()
		log.Error("Data must be JSON format:", incoming, ", But current is: ", err)
		return
	}
	statistics.IncOut()
	e.PushQueue(typex.QueueData{
		In:   nil,
		Out:  e.AllOutEnd()[id],
		E:    e,
		Data: incoming,
	})

}
