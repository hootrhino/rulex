package httpserver

import "github.com/gin-gonic/gin"

// 网络配置结构体
type NetConfig struct {
	DHCP    string   `json:"dhcp"`
	Ip      []string `json:"ip"`
	Gateway string   `json:"gateway"`
	Names   []string `json:"names"`
	Version int      `json:"version"`
}

// 读取Ip状态(静态/动态)  yaml
type T struct {
	Network struct {
		Ethernets struct {
			Eth struct {
				DHCP    string   `yaml:"dhcp4"`
				Ip      []string `yaml:"addresses"`
				Gateway string   `yaml:"gateway4"`
				Names   struct {
					Ip []string `yaml:"addresses"`
				} `yaml:"names"`
			} `yaml:"eth0"`
		} `yaml:"ethernets"`
		Version int `json:"version"`
	} `yaml:"network"`
}

// 读取WIFI状态(静态/动态)  yaml
type WT struct {
	Network struct {
		Ethernets struct {
			Eth struct {
				DHCP    string   `yaml:"dhcp4"`
				Ip      []string `yaml:"addresses"`
				Gateway string   `yaml:"gateway4"`
				Names   struct {
					Ip []string `yaml:"addresses"`
				} `yaml:"names"`
			} `yaml:"wlan0"`
		} `yaml:"ethernets"`
		Version int `json:"version"`
	} `yaml:"network"`
}

// 主要是针对WIFI、时区、IP地址设置

/*
*
* WIFI
*
 */
func SetWifi(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置时间、时区
*
 */
func SetTime(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置静态网络IP等
*
 */
func SetStaticNetwork(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}
