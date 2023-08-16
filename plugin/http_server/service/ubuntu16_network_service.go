package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// # /etc/network/interfaces
//-------------------------------------------
// Static
//-------------------------------------------
// auto lo
// iface lo inet loopback
// auto eth0
// iface eth0 inet static
//     address 192.168.1.100
//     netmask 255.255.255.0
//     gateway 192.168.1.1
//     dns-nameservers 8.8.8.8 8.8.4.4

//-------------------------------------------
// DHCP
//-------------------------------------------
// auto lo
// iface lo inet loopback
// auto eth0
// iface eth0 inet dhcp

type EtcNetworkConfig struct {
	Name        string   `json:"name,omitempty"`
	Interface   string   `json:"interface,omitempty"`
	Address     string   `json:"address,omitempty"`
	Netmask     string   `json:"netmask,omitempty"`
	Gateway     string   `json:"gateway,omitempty"`
	DNS         []string `json:"dns,omitempty"`
	DHCPEnabled bool     `json:"dhcp_enabled,omitempty"`
}

func (nc *EtcNetworkConfig) JsonString() string {
	b, _ := json.Marshal(nc)
	return string(b)
}

/*
*
* 将结构体写入配置文件
sudo systemctl restart networking
sudo service networking restart
*
*/
func ApplyConfig(iface EtcNetworkConfig) error {
	configLines := []string{
		"auto lo",
		"iface lo inet loopback",
		fmt.Sprintf("auto %s", iface.Name),
		fmt.Sprintf("iface %s inet %s", iface.Interface, func(dhcpEnabled bool) string {
			if dhcpEnabled {
				return "dhcp"
			}
			return "static"
		}(iface.DHCPEnabled)),
	}

	if !iface.DHCPEnabled {
		configLines = append(configLines, fmt.Sprintf("    address %s", iface.Address))
		configLines = append(configLines, fmt.Sprintf("    netmask %s", iface.Netmask))
		configLines = append(configLines, fmt.Sprintf("    gateway %s", iface.Gateway))
		configLines = append(configLines, fmt.Sprintf("    dns-nameservers %s\n", strings.Join(iface.DNS, " ")))
	}

	configText := strings.Join(configLines, "\n")
	return os.WriteFile("/etc/network/interfaces", []byte(configText), 0644)
}
