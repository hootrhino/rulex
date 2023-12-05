package dto

type HwPort struct {
	Dhcp4       bool     `yaml:"dhcp4" json:"dhcp4,omitempty"`
	Addresses   []string `yaml:"addresses" json:"addresses,omitempty"`
	Gateway4    string   `yaml:"gateway4" json:"gateway4,omitempty"`
	Nameservers []string `yaml:"nameservers" json:"nameservers,omitempty"`
}
type network struct {
	Version   int               `yaml:"version" json:"version,omitempty"`
	Renderer  string            `yaml:"renderer" json:"renderer,omitempty"`
	Ethernets map[string]HwPort `yaml:"ethernets" json:"ethernets,omitempty"`
}
type NetplanConfigDto struct {
	Network network `yaml:"network" json:"network,omitempty"`
}
