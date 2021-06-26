package x

//
//
//
type outEnd struct {
	Id          string                  `json:"id"`
	Type        string                  `json:"type"`
	State       int                     `json:"state"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Config      *map[string]interface{} `json:"config"`
}

//
//
//
func NewOutEnd(t string,
	n string,
	d string,
	c *map[string]interface{}) *outEnd {
	return &outEnd{
		Id:          MakeUUID("OUTEND"),
		Type:        t,
		State:       0,
		Name:        n,
		Description: d,
		Config:      c,
	}
}
//
//
//
func (out *outEnd) GetConfig(k string) interface{} {
	return (*out.Config)[k]
}
