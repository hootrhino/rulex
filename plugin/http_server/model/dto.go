package model

// NetInterfaceInfo 网络适配器信息
type NetInterfaceInfo struct {
	Name string `json:"name,omitempty"`
	Mac  string `json:"mac,omitempty"`
	Addr string `json:"addr,omitempty"`
}
