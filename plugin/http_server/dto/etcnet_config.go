package service

type EtcNetworkConfig struct {
	Name        string   `json:"name,omitempty"`
	Interface   string   `json:"interface,omitempty"`
	Address     string   `json:"address,omitempty"`
	Netmask     string   `json:"netmask,omitempty"`
	Gateway     string   `json:"gateway,omitempty"`
	DNS         []string `json:"dns,omitempty"`
	DHCPEnabled bool     `json:"dhcp_enabled,omitempty"`
}
