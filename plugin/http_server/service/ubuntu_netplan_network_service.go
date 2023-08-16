package service

import (
	"encoding/json"
	"os"

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
	Dhcp4       bool     `yaml:"dhcp4" json:"dhcp4,omitempty"`
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

/*
*
  - 默认静态IP
    114DNS:
    IPv4: 114.114.114.114, 114.114.115.115
    IPv6: 2400:3200::1, 2400:3200:baba::1
    阿里云DNS:
    IPv4: 223.5.5.5, 223.6.6.6
    腾讯DNS:
    IPv4: 119.29.29.29, 119.28.28.28
    百度DNS:
    IPv4: 180.76.76.76
    DNSPod DNS (也称为Dnspod Public DNS):
    IPv4: 119.29.29.29, 182.254.116.116
*/
func DefaultStaticNetplanConfig() *NetplanConfig {
	return &NetplanConfig{
		Network: Network{
			Version:  2,
			Renderer: "NetworkManager",
			Ethernets: map[string]Interface{
				"eth0": {
					Dhcp4:       false,
					Addresses:   []string{"192.168.128.1/24"},
					Gateway4:    "192.168.128.1",
					Nameservers: []string{"114.114.114.114", "8.8.8.8", "180.76.76.76"},
				},
				"eth1": {
					Dhcp4:       false,
					Addresses:   []string{"192.168.128.1/24"},
					Gateway4:    "192.168.128.2",
					Nameservers: []string{"114.114.114.114", "8.8.8.8", "180.76.76.76"},
				},
			},
		},
	}
}

/*
*
* 默认 DHCP
*
 */
func DefaultDHCPNetplanConfig() *NetplanConfig {
	return &NetplanConfig{
		Network: Network{
			Version:  2,
			Renderer: "NetworkManager",
			Ethernets: map[string]Interface{
				"eth0": {
					Dhcp4: true,
				},
				"eth1": {
					Dhcp4: true,
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

/*
*
* 将配置写入文件并且重启网卡
*
 */
func (nc *NetplanConfig) ApplyConfig() error {
	// sudo netplan apply
	// sudo systemctl restart systemd-networkd
	// sudo service networking restart
	return os.WriteFile("/etc/netplan/001-cfg.yaml", []byte(nc.YAMLString()), 0755)
}
