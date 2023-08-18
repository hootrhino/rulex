package httpserver

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 设置音量
*
 */
func SetVolume(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* WIFI
*
 */
func SetWifi(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		Interface string `json:"interface"`
		SSID      string `json:"ssid"`
		Password  string `json:"password"`
		Security  string `json:"security"` // wpa2-psk wpa3-psk
	}

	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	if !strings.Contains("wpa2-psk wpa3-psk", DtoCfg.Security) {
		c.JSON(common.HTTP_OK,
			common.Error(("Only support 2 valid security algorithm:wpa2-psk,wpa3-psk")))
		return
	}
	if !strings.Contains("wlan0", DtoCfg.Interface) {
		c.JSON(common.HTTP_OK,
			common.Error(("Only support wlan0")))
		return
	}
	UbuntuVersion, err := utils.GetUbuntuVersion()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	NetCfgType := "NETWORK_ETC"
	if (UbuntuVersion == "ubuntu18") ||
		UbuntuVersion == "ubuntu20" ||
		UbuntuVersion == "ubuntu22" ||
		(UbuntuVersion == "ubuntu24") {
		NetCfgType = "NETPLAN"
	}
	MNetCfg := model.MWifiConfig{
		Interface: DtoCfg.Interface,
		SSID:      DtoCfg.SSID,
		Password:  DtoCfg.Password,
		Security:  DtoCfg.Security,
	}
	if err := service.UpdateWlan0Config(MNetCfg); err != nil {
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if NetCfgType == "NETPLAN" { // Ubuntu18极其以后
		ApplyNewestNetplanWlanConfig()
	}
	// 保存到数据库, 并且写入配置
	c.JSON(common.HTTP_OK, common.OkWithData(DtoCfg))

}

/*
*
  - 设置时间、时区
  - sudo date -s "2023-08-07 15:30:00"
    获取时间: date "+%Y-%m-%d %H:%M:%S" -> 2023-08-07 15:30:00
*/
func SetTime(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置静态网络IP等, 当前只支持Linux 其他的没测试暂时不予支持

	{
	  "name": "eth0",
	  "interface": "eth0",
	  "address": "192.168.1.100",
	  "netmask": "255.255.255.0",
	  "gateway": "192.168.1.1",
	  "dns": ["8.8.8.8", "8.8.4.4"],
	  "dhcp_enabled": false
	}
*/
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func SetEthNetwork(c *gin.Context, hh *HttpApiServer) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("Set Static Network Not Support:"+runtime.GOOS))
		return
	}
	type Form struct {
		Interface   string   `json:"interface"` // eth1 eth0
		Address     string   `json:"address"`
		Netmask     string   `json:"netmask"`
		Gateway     string   `json:"gateway"`
		DNS         []string `json:"dns"`
		DHCPEnabled bool     `json:"dhcp_enabled"`
	}

	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	if !strings.Contains("eth1 eth0", DtoCfg.Interface) {
		c.JSON(common.HTTP_OK,
			common.Error(("Only have 2 valid interface:eth1 and eth0")))
		return
	}
	if !isValidIP(DtoCfg.Address) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid IP:" + DtoCfg.Address)))
		return
	}
	if !isValidIP(DtoCfg.Gateway) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid Gateway IP:" + DtoCfg.Address)))
		return
	}
	if !isValidSubnetMask(DtoCfg.Netmask) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid SubnetMask:" + DtoCfg.Address)))
		return
	}
	for _, dns := range DtoCfg.DNS {
		if !isValidIP(dns) {
			c.JSON(common.HTTP_OK,
				common.Error(("Invalid DNS IP:" + DtoCfg.Address)))
			return
		}
	}
	UbuntuVersion, err := utils.GetUbuntuVersion()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	NetCfgType := "NETWORK_ETC"
	if (UbuntuVersion == "ubuntu18") ||
		UbuntuVersion == "ubuntu20" ||
		UbuntuVersion == "ubuntu22" ||
		(UbuntuVersion == "ubuntu24") {
		NetCfgType = "NETPLAN"
	}
	MNetCfg := model.MNetworkConfig{
		Type:        NetCfgType,
		Interface:   DtoCfg.Interface,
		Address:     DtoCfg.Address,
		Netmask:     DtoCfg.Netmask,
		Gateway:     DtoCfg.Gateway,
		DNS:         DtoCfg.DNS,
		DHCPEnabled: DtoCfg.DHCPEnabled,
	}
	if DtoCfg.Interface == "eth0" {
		if err := service.UpdateEth0Config(MNetCfg); err != nil {
			if err != nil {
				c.JSON(common.HTTP_OK, common.Error400(err))
				return
			}
		}
	}
	if DtoCfg.Interface == "eth1" {
		if err := service.UpdateEth1Config(MNetCfg); err != nil {
			if err != nil {
				c.JSON(common.HTTP_OK, common.Error400(err))
				return
			}
		}
	}

	if NetCfgType == "NETWORK_ETC" { // Ubuntu16
		ApplyNewestEtcConfig()
	}
	if NetCfgType == "NETPLAN" { // Ubuntu18极其以后
		ApplyNewestNetplanEthConfig()
	}
	// 保存到数据库, 并且写入配置
	c.JSON(common.HTTP_OK, common.OkWithData(DtoCfg))

}

/*
*
* 生成最新的无线配置
*
 */
func ApplyNewestNetplanWlanConfig() error {
	MWlan0, err := service.GetWlan0Config()
	if err != nil {
		return err
	}
	Wlan0Config := service.WlanConfig{
		Wlan0: service.WLANInterface{
			Interface: MWlan0.Interface,
			SSID:      MWlan0.SSID,
			Password:  MWlan0.Password,
			Security:  MWlan0.Security,
		},
	}
	fmt.Println(Wlan0Config.YAMLString())
	return nil
}

/*
*
* 生成YAML
*
 */
func ApplyNewestNetplanEthConfig() error {
	Eth0, err := service.GetEth0Config()
	if err != nil {
		return err
	}
	Eth1, err := service.GetEth1Config()
	if err != nil {
		return err
	}

	NetplanConfig := service.NetplanConfig{
		Network: service.Network{
			Version:  2,
			Renderer: "NetworkManager",
			Ethernets: service.EthInterface{
				Eth0: service.HwInterface{
					Dhcp4:       Eth0.DHCPEnabled,
					Addresses:   []string{Eth0.Address},
					Gateway4:    Eth0.Gateway,
					Nameservers: Eth0.DNS,
				},
				Eth1: service.HwInterface{
					Dhcp4:       Eth1.DHCPEnabled,
					Addresses:   []string{Eth1.Address},
					Gateway4:    Eth1.Gateway,
					Nameservers: Eth1.DNS,
				},
			},
		},
	}
	fmt.Println(NetplanConfig.YAMLString())
	return nil

}

/*
*
* 生成 etc 文件
*
 */
func ApplyNewestEtcConfig() error {
	// 取出最新的配置
	MNetworkConfigs, err := service.GetAllNetConfig()
	if err != nil {
		return err
	}
	etcFileContent := ""
	for _, MNetworkConfig := range MNetworkConfigs {
		NetworkConfig := service.EtcNetworkConfig{
			Interface:   MNetworkConfig.Interface,
			Address:     MNetworkConfig.Address,
			Netmask:     MNetworkConfig.Netmask,
			Gateway:     MNetworkConfig.Gateway,
			DNS:         MNetworkConfig.DNS,
			DHCPEnabled: MNetworkConfig.DHCPEnabled,
		}
		etcFileContent += NetworkConfig.GenEtcConfig() + "\n"
	}
	fmt.Println(etcFileContent)
	return nil
}
func isValidSubnetMask(mask string) bool {
	// 分割子网掩码为4个整数
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return false
	}

	// 将每个部分转换为整数
	var octets [4]int
	for i, part := range parts {
		octet, err := strconv.Atoi(part)
		if err != nil || octet < 0 || octet > 255 {
			return false
		}
		octets[i] = octet
	}

	// 判断是否为有效的子网掩码
	var bits int
	for _, octet := range octets {
		bits += bitsInByte(octet)
	}

	return bits >= 1 && bits <= 32
}

func bitsInByte(b int) int {
	count := 0
	for b > 0 {
		count += b & 1
		b >>= 1
	}
	return count
}
