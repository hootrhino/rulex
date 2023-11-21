package dto

// {
// 	"name": "eth0",
// 	"interface": "eth0",
// 	"address": "192.168.1.100",
// 	"netmask": "255.255.255.0",
// 	"gateway": "192.168.1.1",
// 	"dns": ["8.8.8.8", "8.8.4.4"],
// 	"dhcp_enabled": false
// }

type EtcNetworkConfig struct {
	Name        string   `json:"name"`      // eth1 eth0
	Interface   string   `json:"interface"` // eth1 eth0
	Address     string   `json:"address"`
	Netmask     string   `json:"netmask"`
	Gateway     string   `json:"gateway"`
	DNS         []string `json:"dns"`
	DHCPEnabled bool     `json:"dhcp_enabled"`
}
