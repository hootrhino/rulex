package x

import "sync"

//
type inEnd struct {
	sync.Mutex
	Id          string                  `json:"id"`
	State       State                   `json:"state"`
	Type        string                  `json:"type"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Binds       *map[string]rule        `json:"-"`
	Config      *map[string]interface{} `json:"config"`
}

func (in *inEnd) GetState() State {
	in.Lock()
	defer in.Unlock()
	return in.State
}

//
func (in *inEnd) SetState(s State) {
	in.Lock()
	defer in.Unlock()
	in.State = s
}

//
func NewInEnd(t string,
	n string,
	d string,
	c *map[string]interface{}) *inEnd {

	return &inEnd{
		Id:          MakeUUID("INEND"),
		Type:        t,
		Name:        n,
		Description: d,
		Binds:       &map[string]rule{},
		Config:      c,
	}
}
func (in *inEnd) GetConfig(k string) interface{} {
	return (*in.Config)[k]
}
