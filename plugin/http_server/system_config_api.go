package httpserver

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/glogger"
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
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	type Form struct {
		Volume int `json:"volume"`
	}
	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	v, err := service.SetVolume(DtoCfg.Volume)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(v))

}

/*
*
* 获取音量的值
*
 */
func GetVolume(c *gin.Context, hh *HttpApiServer) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	v, err := service.GetVolume()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if v == "" {
		c.JSON(common.HTTP_OK, common.Error("Volume get failed, please check system"))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"volume": v,
	}))
}

/*
*
* WIFI
*
 */
func GetWifi(c *gin.Context, hh *HttpApiServer) {
	MWifiConfig, err := service.GetWlan0Config()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Cfg := service.WlanConfig{
		Wlan0: service.WLANInterface{
			Interface: MWifiConfig.Interface,
			SSID:      MWifiConfig.SSID,
			Password:  MWifiConfig.Password,
			Security:  MWifiConfig.Security,
		},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Cfg))

}

/*
*
*
*通过nmcli配置WIFI
 */
func SetWifi(c *gin.Context, hh *HttpApiServer) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
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
	if !utils.SContains([]string{"wpa2-psk", "wpa3-psk"}, DtoCfg.Security) {
		c.JSON(common.HTTP_OK,
			common.Error(("Only support 2 valid security algorithm:wpa2-psk,wpa3-psk")))
		return
	}
	if !utils.SContains([]string{"wlan0"}, DtoCfg.Interface) {
		c.JSON(common.HTTP_OK, common.Error(("Only support wlan0")))
		return
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
	/*
	*
	* 全部使用nmcli操作
	*
	 */
	ApplyNewestEtcWlanConfig()
	service.EtcApply()

	// 保存到数据库, 并且写入配置
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 生成最新的ETC配置
*
 */
func ApplyNewestEtcEthConfig() error {
	MEth0, err := service.GetEth0Config()
	if err != nil {
		return err
	}
	MEth1, err := service.GetEth1Config()
	if err != nil {
		return err
	}
	EtcEth0Cfg := service.EtcNetworkConfig{
		Interface:   MEth0.Interface,
		Address:     MEth0.Address,
		Netmask:     MEth0.Netmask,
		Gateway:     MEth0.Gateway,
		DNS:         MEth0.DNS,
		DHCPEnabled: *MEth0.DHCPEnabled,
	}
	EtcEth1Cfg := service.EtcNetworkConfig{
		Interface:   MEth1.Interface,
		Address:     MEth1.Address,
		Netmask:     MEth1.Netmask,
		Gateway:     MEth1.Gateway,
		DNS:         MEth1.DNS,
		DHCPEnabled: *MEth1.DHCPEnabled,
	}
	loopBack := "# DON'T EDIT THIS FILE!\nauto lo\niface lo inet loopback\n"
	return os.WriteFile("/etc/network/interfaces",
		[]byte(
			loopBack+
				EtcEth0Cfg.GenEtcConfig()+
				"\n"+
				EtcEth1Cfg.GenEtcConfig()+"\n"), 0755)

}

/*
*
  - 设置时间、时区
  - sudo date -s "2023-08-07 15:30:00"
    获取时间: date "+%Y-%m-%d %H:%M:%S" -> 2023-08-07 15:30:00
*/
func SetSystemTime(c *gin.Context, hh *HttpApiServer) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	type Form struct {
		Time string `json:"time"`
	}
	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}

	err := service.SetSystemTime(DtoCfg.Time)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 用来测试生成各种网络配置
*
 */
func TestGenEtcNetCfg(c *gin.Context, hh *HttpApiServer) {

	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 时区设置
*
 */
func SetSystemTimeZone(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		Timezone string `json:"timezone"`
	}
	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}

	if !validTimeZone(DtoCfg.Timezone) {
		c.JSON(common.HTTP_OK, common.Error("Invalid timezone:"+DtoCfg.Timezone))

	}
	if err := service.SetTimeZone(DtoCfg.Timezone); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	}
	c.JSON(common.HTTP_OK, common.Ok())

}
func GetSystemTimeZone(c *gin.Context, hh *HttpApiServer) {
	c.JSON(common.HTTP_OK, common.OkWithData(service.GetTimeZone()))

}

/*
*
* 获取系统时间
*
 */
func GetSystemTime(c *gin.Context, hh *HttpApiServer) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	SysTime := service.GetSystemTime()
	c.JSON(common.HTTP_OK, common.OkWithData(SysTime))
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
func validTimeZone(timezone string) bool {
	// 使用正则表达式来匹配时区格式
	// 时区格式应该类似于 "America/New_York" 或 "Asia/Shanghai"
	// 这里使用了简单的正则表达式，你可以根据需要进行调整
	regexPattern := `^[A-Za-z]+/[A-Za-z_]+$`
	regex := regexp.MustCompile(regexPattern)

	return regex.MatchString(timezone)
}

/*
*
* 展示网络配置信息
*
 */
func GetEthNetwork(c *gin.Context, hh *HttpApiServer) {
	MEth0, err := service.GetEth0Config()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return

	}
	MEth1, err := service.GetEth1Config()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Eth0Cfg := service.EtcNetworkConfig{
		Interface:   MEth0.Interface,
		Address:     MEth0.Address,
		Netmask:     MEth0.Address,
		Gateway:     MEth0.Address,
		DNS:         MEth0.DNS,
		DHCPEnabled: *MEth0.DHCPEnabled,
	}
	Eth1Cfg := service.EtcNetworkConfig{
		Interface:   MEth1.Interface,
		Address:     MEth1.Address,
		Netmask:     MEth1.Address,
		Gateway:     MEth1.Address,
		DNS:         MEth1.DNS,
		DHCPEnabled: *MEth1.DHCPEnabled,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]service.EtcNetworkConfig{
		"eth0": Eth0Cfg,
		"eth1": Eth1Cfg,
	}))

}
func SetEthNetwork(c *gin.Context, hh *HttpApiServer) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
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
	if !utils.SContains([]string{"eth1", "eth0"}, DtoCfg.Interface) {
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
		DHCPEnabled: &DtoCfg.DHCPEnabled,
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
	/*
	*
	* 全部采用nmcli
	*
	 */
	ApplyNewestEtcEthConfig()
	service.EtcApply()
	c.JSON(common.HTTP_OK, common.Ok())

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
	// fmt.Println(Wlan0Config.YAMLString())
	return Wlan0Config.ApplyWlan0Config()
}

/*
*
* ubuntu1604网络, 使用一个 nmcli 指令
*
 */
func ApplyNewestEtcWlanConfig() error {
	MWlan0, err := service.GetWlan0Config()
	if err != nil {
		return err
	}
	// nmcli dev wifi connect SSID password pwd
	s := "nmcli dev wifi connect \"%s\" password \"%s\""
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf(s, MWlan0.SSID, MWlan0.Password))
	out, err := cmd.CombinedOutput()
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	glogger.GLogger.Info(out)
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
					Addresses:   []string{Eth0.Address + "/24"},
					Gateway4:    Eth0.Gateway,
					Nameservers: Eth0.DNS,
				},
				Eth1: service.HwInterface{
					Dhcp4:       Eth1.DHCPEnabled,
					Addresses:   []string{Eth1.Address + "/24"},
					Gateway4:    Eth1.Gateway,
					Nameservers: Eth1.DNS,
				},
			},
		},
	}
	// fmt.Println(NetplanConfig.YAMLString())
	return NetplanConfig.ApplyEthConfig()

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
