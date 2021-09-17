package typex

import (
	"rulex/utils"
	"sync"
)

//
//
//
type OutEnd struct {
	sync.Mutex
	Id          string                  `json:"id"`
	Type        TargetType              `json:"type"`
	State       ResourceState           `json:"state"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Config      *map[string]interface{} `json:"config"`
	Target      XTarget                 `json:"-"`
}

func (o *OutEnd) GetState() ResourceState {
	o.Lock()
	defer o.Unlock()
	return o.State
}

//
func (o *OutEnd) SetState(s ResourceState) {
	o.Lock()
	defer o.Unlock()
	o.State = s
}

//
//
//
func NewOutEnd(t string,
	n string,
	d string,
	c *map[string]interface{}) *OutEnd {
	return &OutEnd{
		Id:          utils.MakeUUID("OUTEND"),
		Type:        t,
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
	return (*out.Config)[k]
}
