package test

type Mapping struct {
	Type  string
	Value string
}
type Spec []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	DataType struct {
		Type    string  `json:"type"`
		Mapping Mapping `json:"mapping"`
	} `json:"dataType"`
}
type Define struct {
	Type  string `json:"type"`
	Specs []Spec `json:"specs"`
}
type Propertie struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Mode     string `json:"mode"`
	Define   Define `json:"define"`
	Required bool   `json:"required"`
}
type Param struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Define Define `json:"define"`
}
type Event struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Type     string  `json:"type"`
	Params   []Param `json:"params"`
	Required bool    `json:"required"`
}
type InDefine struct {
	Type    string  `json:"type"`
	Mapping Mapping `json:"mapping"`
}
type Input struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Define InDefine `json:"define"`
}
type OutDefine struct {
	Type    string  `json:"type"`
	Mapping Mapping `json:"mapping"`
}
type Output struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Define OutDefine `json:"define"`
}
type Action struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Input    []Input  `json:"input"`
	Output   []Output `json:"output"`
	Required bool     `json:"required"`
}
type Profile struct {
	ProductID  string `json:"ProductId"`
	CategoryID string `json:"CategoryId"`
}
type Schema struct {
	Version    string      `json:"version"`
	Properties []Propertie `json:"properties"`
	Events     []Event     `json:"events"`
	Actions    []Action    `json:"actions"`
	Profile    Profile     `json:"profile"`
}
