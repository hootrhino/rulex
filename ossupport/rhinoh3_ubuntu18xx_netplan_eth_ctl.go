package ossupport

import (
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
type HwPort struct {
	Dhcp4       *bool    `yaml:"dhcp4" json:"dhcp4"`
	Addresses   []string `yaml:"addresses" json:"addresses"`
	Gateway4    string   `yaml:"gateway4" json:"gateway4"`
	Nameservers []string `yaml:"nameservers" json:"nameservers"`
}

type EthInterface struct {
	Eth0 HwPort `yaml:"eth0" json:"eth0"`
	Eth1 HwPort `yaml:"eth1" json:"eth1"`
}
type Network struct {
	Version   int          `yaml:"version" json:"version"`
	Renderer  string       `yaml:"renderer" json:"renderer"`
	Ethernets EthInterface `yaml:"ethernets" json:"ethernets"`
}
type NetplanConfig struct {
	Network Network `yaml:"network" json:"network"`
}

/*
*
* 默认 DHCP
*
 */

func (nc *NetplanConfig) FromYaml(jsons string) error {
	return yaml.Unmarshal([]byte(jsons), nc)
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
func (nc *NetplanConfig) ApplyEthConfig() error {
	// sudo netplan apply
	// sudo systemctl restart systemd-networkd
	// sudo ossupport networking restart
	// return os.WriteFile("/etc/netplan/001-eth.yaml", []byte(nc.YAMLString()), 0755)
	return os.WriteFile("/etc/netplan/001-eth.yaml", []byte(nc.YAMLString()), 0755)
}
