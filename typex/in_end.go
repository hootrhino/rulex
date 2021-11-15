package typex

import (
	"rulex/utils"
	"sync"
)

//
type InEnd struct {
	sync.Mutex
	//
	UUID        string          `json:"uuid"`
	State       ResourceState   `json:"state"`
	Type        InEndType       `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Binds       map[string]Rule `json:"-"`
	//
	Config   map[string]interface{} `json:"config"`
	Resource XResource              `json:"-"`
}

func (in *InEnd) GetState() ResourceState {
	in.Lock()
	defer in.Unlock()
	return in.State
}

//
func (in *InEnd) SetState(s ResourceState) {
	in.Lock()
	defer in.Unlock()
	in.State = s
}

//
func NewInEnd(t string,
	n string,
	d string,
	c map[string]interface{}) *InEnd {

	return &InEnd{
		UUID:        utils.MakeUUID("INEND"),
		Type:        InEndType(t),
		Name:        n,
		Description: d,
		Binds:       map[string]Rule{},
		Config:      c,
	}
}
func (in *InEnd) GetConfig(k string) interface{} {
	return (in.Config)[k]
}
