package typex

import (
	"rulex/utils"
)

//
//
//
type OutEnd struct {
	UUID        string        `json:"uuid"`
	State       ResourceState `json:"state"`
	Type        TargetType    `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	//
	Config map[string]interface{} `json:"config"`
	Target XTarget                `json:"-"`
}

func (o *OutEnd) GetState() ResourceState {
	return o.State
}

//
func (o *OutEnd) SetState(s ResourceState) {
	o.State = s
}

//
//
//
func NewOutEnd(t string,
	n string,
	d string,
	c map[string]interface{}) *OutEnd {
	return &OutEnd{
		UUID:        utils.MakeUUID("OUTEND"),
		Type:        TargetType(t),
		State:       DOWN,
		Name:        n,
		Description: d,
		Config:      c,
	}
}

//
//
//
func (out *OutEnd) GetConfig(k string) interface{} {
	return (out.Config)[k]
}
