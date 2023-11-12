package ossupport

import (
	"fmt"
	"os"
)

// network:
//
//	version: 2
//	renderer: networkd
//	wifis:
//	  wlan0:
//	    dhcp4: yes
//	    access-points:
//	      "YourWiFiSSID":
//	        password: "YourWiFiPassword"

type WLANInterface struct {
	Interface string `yaml:"-" json:"interface"`
	SSID      string `yaml:"-" json:"ssid"`
	Password  string `yaml:"-" json:"password"`
	Security  string `yaml:"-" json:"security"`
}

type WlanConfig struct {
	Wlan0 WLANInterface `yaml:"-" json:"wlan0"`
}

/*
*
* 专门配置WIFI
*
 */
func (nc *WlanConfig) YAMLString() string {
	var netplanWLAN0Config string = fmt.Sprintf(
		`network:
  version: 2
  renderer: NetworkManager
  wifis:
    wlan0:
      dhcp4: yes
      access-points:
        "%s":
          password: "%s"
`, nc.Wlan0.SSID, nc.Wlan0.Password)
	return netplanWLAN0Config
}

/*
*
* 将配置写入文件并且重启网卡
*
 */
func (nc *WlanConfig) ApplyWlan0Config() error {
	return os.WriteFile("/etc/netplan/001-wlan.yaml", []byte(nc.YAMLString()), 0755)
}
