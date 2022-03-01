package typex

import "rulex/utils"

//
type InEnd struct {
	//
	UUID        string          `json:"uuid"`
	State       SourceState     `json:"state"`
	Type        InEndType       `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Binds       map[string]Rule `json:"-"`
	//
	Config        map[string]interface{} `json:"config"`
	DataModelsMap map[string]XDataModel  `json:"-"`
	Source        XSource                `json:"-"`
}

func (in *InEnd) GetState() SourceState {
	return in.State
}

//
func (in *InEnd) SetState(s SourceState) {
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
