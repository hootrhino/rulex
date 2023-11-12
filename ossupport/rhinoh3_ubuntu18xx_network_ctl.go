package ossupport

import (
	"encoding/json"
	"fmt"
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
	Interface   string   `json:"interface"`
	Address     string   `json:"address"`
	Netmask     string   `json:"netmask"`
	Gateway     string   `json:"gateway"`
	DNS         []string `json:"dns"`
	DHCPEnabled bool     `json:"dhcp_enabled"`
}

func (nc *EtcNetworkConfig) JsonString() string {
	b, _ := json.Marshal(nc)
	return string(b)
}

/*
*
* 将结构体写入配置文件
sudo systemctl restart networking
sudo ossupport networking restart
*
*/
func (iface *EtcNetworkConfig) GenEtcConfig() string {
	configLines := []string{
		fmt.Sprintf("auto %s", iface.Interface),
		fmt.Sprintf("iface %s inet %s", iface.Interface, func(dhcpEnabled bool) string {
			if dhcpEnabled {
				return "dhcp"
			}
			return "static"
		}(iface.DHCPEnabled)),
	}
	if iface.DHCPEnabled {
		return strings.Join(configLines, "\n")
	}
	configLines = append(configLines, fmt.Sprintf("    address %s", iface.Address))
	configLines = append(configLines, fmt.Sprintf("    netmask %s", iface.Netmask))
	configLines = append(configLines, fmt.Sprintf("    gateway %s", iface.Gateway))
	configLines = append(configLines, fmt.Sprintf("    dns-nameservers %s\n", strings.Join(iface.DNS, " ")))
	configText := strings.Join(configLines, "\n")
	return configText
}
