package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

/*
* Ubuntu 18 以后的版本才支持
/etc/netplan/01-netcfg.yaml
network:
  version: 2
  renderer: networkd
  ethernets:
    enp0s9:
      dhcp4: no
      addresses:
        - 192.168.121.221/24
      gateway4: 192.168.121.1
      nameservers:
          addresses: [8.8.8.8, 1.1.1.1]
*
*/
//
// 读取Ip状态(静态/动态)  yaml
type Interface struct {
	Dhcp4       string   `yaml:"dhcp4" json:"dhcp4,omitempty"`
	Addresses   []string `yaml:"addresses" json:"addresses,omitempty"`
	Gateway4    string   `yaml:"gateway4" json:"gateway4,omitempty"`
	Nameservers []string `yaml:"nameservers" json:"nameservers,omitempty"`
}
type Network struct {
	Version   int                  `yaml:"version" json:"version,omitempty"`
	Renderer  string               `yaml:"renderer" json:"renderer,omitempty"`
	Ethernets map[string]Interface `yaml:"ethernets" json:"ethernets,omitempty"`
}
type NetplanConfig struct {
	Network Network `yaml:"network" json:"network,omitempty"`
}

func NewNetplanConfig() *NetplanConfig {
	return &NetplanConfig{
		Network: Network{
			Version:  2,
			Renderer: "NetworkManager",
			Ethernets: map[string]Interface{
				"eth0": {
					Dhcp4:       "no",
					Addresses:   []string{"192.168.128.1/24"},
					Gateway4:    "192.168.128.1",
					Nameservers: []string{"114.114.114.114"},
				},
				"eth1": {
					Dhcp4:       "no",
					Addresses:   []string{"192.168.128.1/24"},
					Gateway4:    "192.168.128.2",
					Nameservers: []string{"114.114.114.114"},
				},
			},
		},
	}
}
func (nc *NetplanConfig) FromJson(jsons string) error {
	return json.Unmarshal([]byte(jsons), nc)
}

func (nc *NetplanConfig) FromYaml(jsons string) error {
	return yaml.Unmarshal([]byte(jsons), nc)
}

func (nc *NetplanConfig) JsonString() string {
	b, _ := json.Marshal(nc)
	return string(b)
}
func (nc *NetplanConfig) YAMLString() string {
	b, _ := yaml.Marshal(nc)
	return string(b)
}

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
* 解析配置文件
*
 */
func ParseEtcFile(content string) []EtcNetworkConfig {
	lines := strings.Split(content, "\n")

	var interfaces []EtcNetworkConfig
	var currentInterface EtcNetworkConfig

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "auto":
			currentInterface = EtcNetworkConfig{Name: fields[1]}
		case "iface":
			if len(fields) < 3 {
				continue
			}
			currentInterface.Interface = fields[2]
		case "address":
			if len(fields) < 2 {
				continue
			}
			currentInterface.Address = fields[1]
		case "netmask":
			if len(fields) < 2 {
				continue
			}
			currentInterface.Netmask = fields[1]
		case "gateway":
			if len(fields) < 2 {
				continue
			}
			currentInterface.Gateway = fields[1]
		case "dns-nameservers":
			if len(fields) < 2 {
				continue
			}
			currentInterface.DNS = fields[1:]
		case "dhcp":
			if len(fields) > 1 && fields[1] == "dhcp" {
				currentInterface.DHCPEnabled = true
			}
		}

		if len(currentInterface.Interface) > 0 {
			interfaces = append(interfaces, currentInterface)
		}
	}

	return interfaces
}

/*
*
* 将结构体写入配置文件
*
 */
func WriteInterfaceConfig(filePath string, iface EtcNetworkConfig) error {
	configLines := []string{
		"auto lo",
		"iface lo inet loopback",
		fmt.Sprintf("auto %s", iface.Name),
		fmt.Sprintf("iface %s inet %s", iface.Interface, getInetType(iface.DHCPEnabled)),
	}

	if !iface.DHCPEnabled {
		configLines = append(configLines, fmt.Sprintf("    address %s", iface.Address))
		configLines = append(configLines, fmt.Sprintf("    netmask %s", iface.Netmask))
		configLines = append(configLines, fmt.Sprintf("    gateway %s", iface.Gateway))
		configLines = append(configLines, fmt.Sprintf("    dns-nameservers %s", strings.Join(iface.DNS, " ")))
	}

	configText := strings.Join(configLines, "\n")
	return os.WriteFile(filePath, []byte(configText), 0644)
}
func getInetType(dhcpEnabled bool) string {
	if dhcpEnabled {
		return "dhcp"
	}
	return "static"
}
