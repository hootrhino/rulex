package typex

import "github.com/hootrhino/rulex/utils"

type InEnd struct {
	//
	UUID        string          `json:"uuid"`
	State       SourceState     `json:"state"`
	Type        InEndType       `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	BindRules   map[string]Rule `json:"-"`
	//
	Config        map[string]interface{} `json:"config"`
	DataModelsMap map[string]XDataModel  `json:"-"`
	Source        XSource                `json:"-"`
}


func NewInEnd(Type InEndType,
	n string,
	d string,
	c map[string]interface{}) *InEnd {

	return &InEnd{
		UUID:        utils.InUuid(),
		Type:        Type,
		Name:        n,
		Description: d,
		BindRules:   map[string]Rule{},
		Config:      c,
	}
}
func (in *InEnd) GetConfig(k string) interface{} {
	return (in.Config)[k]
}
