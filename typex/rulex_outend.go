package typex

import (
	"github.com/hootrhino/rulex/utils"
)

type OutEnd struct {
	UUID        string      `json:"uuid"`
	State       SourceState `json:"state"`
	Type        TargetType  `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	//
	Config map[string]interface{} `json:"config"`
	Target XTarget                `json:"-"`
}


func NewOutEnd(t TargetType,
	n string,
	d string,
	c map[string]interface{}) *OutEnd {
	return &OutEnd{
		UUID:        utils.MakeUUID("OUTEND"),
		Type:        TargetType(t),
		State:       SOURCE_DOWN,
		Name:        n,
		Description: d,
		Config:      c,
	}
}

func (out *OutEnd) GetConfig(k string) interface{} {
	return (out.Config)[k]
}
