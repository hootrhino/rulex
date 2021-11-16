package stdlib

import (
	"rulex/statistics"
	"rulex/typex"
)

func handleDataFormat(e typex.RuleX, id string, incoming string) {
	// data := &map[string]interface{}{}
	// err := json.Unmarshal([]byte(incoming), data)
	// if err != nil {
	// 	statistics.IncOutFailed()
	// 	log.Error("Data must be JSON format:", incoming, ", But current is: ", err)
	// 	return
	// }
	statistics.IncOut()
	outEnds := (e.AllInEnd())
	outEnd, _ := outEnds.Load(id)
	e.PushQueue(typex.QueueData{
		In:   nil,
		Out:  outEnd.(*typex.OutEnd),
		E:    e,
		Data: incoming,
	})

}
