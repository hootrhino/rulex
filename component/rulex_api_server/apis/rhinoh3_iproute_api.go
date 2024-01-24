// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package apis

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	leases "github.com/hootrhino/go-dhcpd-leases"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/ossupport"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* DHCP 配置
*
 */
type DHCPVo struct {
	Iface       string `json:"iface" validate:"required"`         // 用来做子网的那个网卡的网卡名
	Ip          string `json:"ip" validate:"required"`            // 用来做子网的那个网卡的IP地址
	Gateway     string `json:"gateway" validate:"required"`       // 用来做子网的那个网卡的网关
	Network     string `json:"network" validate:"required"`       // 用来做子网的那个网卡的网段
	Netmask     string `json:"netmask" validate:"required"`       // 用来做子网的那个网卡子网掩码
	IpPoolBegin string `json:"ip_pool_begin" validate:"required"` // DHCP IP地址池起始
	IpPoolEnd   string `json:"ip_pool_end" validate:"required"`   // DHCP IP地址池结束
	//------------------------------------
	// IP 路由方向, 默认 ETH1 透传到 4G
	//------------------------------------
	IfaceFrom string `json:"iface_from" validate:"required"` // 流量入口,固定ETH1
	IfaceTo   string `json:"iface_to" validate:"required"`   // 流量出口,固定4G
}

func SetDHCP(c *gin.Context, ruleEngine typex.RuleX) {
	dhcpVo := DHCPVo{}
	if err := c.ShouldBindJSON(&dhcpVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	err1 := service.UpdateIpRoute(model.MIpRoute{
		Iface:       dhcpVo.Iface,
		Ip:          dhcpVo.Ip,
		Gateway:     dhcpVo.Gateway,
		Network:     dhcpVo.Network,
		Netmask:     dhcpVo.Netmask,
		IpPoolBegin: dhcpVo.IpPoolBegin,
		IpPoolEnd:   dhcpVo.IpPoolEnd,
		IfaceFrom:   dhcpVo.IfaceFrom,
		IfaceTo:     dhcpVo.IfaceTo,
	})
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(200, common.Ok())
}
func GetDHCP(c *gin.Context, ruleEngine typex.RuleX) {
	Model, err := service.GetDefaultIpRoute()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(200, common.OkWithData(DHCPVo{
		Iface:       Model.Iface,
		Ip:          Model.Ip,
		Gateway:     Model.Gateway,
		Network:     Model.Network,
		Netmask:     Model.Netmask,
		IpPoolBegin: Model.IpPoolBegin,
		IpPoolEnd:   Model.IpPoolEnd,
		IfaceFrom:   Model.IfaceFrom,
		IfaceTo:     Model.IfaceTo,
	}))
}

/*
*
  - 更新默认路由
    1 取上一次的路由
    2 配置最新的
    4 删除上一次的
    5 更新最新的路由

*
*/
type IpRouteVo struct {
	UUID  string `json:"uuid"`
	Ip    string `json:"ip" validate:"required"`
	Iface string `json:"iface" validate:"required"`
}

/*
*
* 获取上一次的路由
*
 */
func GetOldDefaultIpRoute(c *gin.Context, ruleEngine typex.RuleX) {
	MIpRoute, err := service.IpRouteDetail()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(IpRouteVo{
		Ip:    MIpRoute.Ip,
		Iface: MIpRoute.Iface,
	}))

}

/*
*
* 设置默认路由
*
 */

func SetNewDefaultIpRoute(c *gin.Context, ruleEngine typex.RuleX) {
	form := IpRouteVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ifaces, err1 := ossupport.GetBSPNetIfaces()
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	if !utils.SContains(ifaces, form.Iface) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Only support iface:%v", ifaces)))
		return
	}
	err3 := service.UpdateIpRoute(model.MIpRoute{
		Ip:    form.Ip,
		Iface: form.Iface,
	})
	if err3 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err3))
		return
	}
	err2 := service.UpdateDefaultRoute(form.Ip, form.Iface)
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取当前在线的DHCP主机列表
*
 */
type DhcpLeaseVo struct {
	MacAddress string `json:"mac_address"` // MAC地址
	IpAddress  string `json:"ip_address"`  // IP地址
	Hostname   string `json:"hostname"`    // 主机名
}

/*
*
* 解析/var/lib/dhcp/dhcpd.leases文件获取DHCP客户端
*
 */
func GetDhcpClients(c *gin.Context, ruleEngine typex.RuleX) {
	f, err := os.Open("/var/lib/dhcp/dhcpd.leases")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	}
	leases := leases.Parse(f)
	Clients := []DhcpLeaseVo{}
	for _, lease := range leases {
		Clients = append(Clients, DhcpLeaseVo{
			IpAddress: lease.IP.String(),
			Hostname: func() string {
				if lease.ClientHostname != "" {
					return lease.ClientHostname
				}
				return "UNKNOWN-HOSTNAME"
			}(),
			MacAddress: lease.Hardware.MAC,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Clients))
}

/*
*
* 删除某一个DHCP客户端
*
 */
func DeleteDhcpClient(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 清空DHCP表
*
 */
func CleanDhcpClients(c *gin.Context, ruleEngine typex.RuleX) {
	f, err := os.Open("/var/lib/dhcp/dhcpd.leases")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	}
	_, err1 := f.Write(nil)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))

	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
*nmcli device status

	DEVICE           TYPE      STATE      CONNECTION
	eth0             ethernet  connected  eth0
	usb0             ethernet  connected  usb0
	wlx0cc6551c5026  wifi      connected  iotlab4072
	eth1             ethernet  connected  eth1
	lo               loopback  unmanaged  --
*/
type networkDevice struct {
	// 网卡名称
	Device string `json:"device"`
	// 网卡类型
	// ethernet：以太网
	// wifi：WiFi
	// bridge：桥接设备
	Type string `json:"type"`
	// 网络状态
	// connected：已连接到。
	// disconnected：未连接。
	// unmanaged：系统默认
	// unavailable：网络不可用。
	State string `json:"state"`
	// 网络名称
	Connection string `json:"connection"`
}

func GetNmcliDeviceStatus(c *gin.Context, ruleEngine typex.RuleX) {

	cmd := exec.Command("nmcli", "device", "status")
	output, err := cmd.Output()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	devices, err := parseNmcliDeviceStatus(string(output))
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(devices))
}

// parseNetworkDevices 解析网络设备信息
func parseNmcliDeviceStatus(input string) ([]networkDevice, error) {
	var devices []networkDevice

	// 将输入按行分割
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "DEVICE") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		device := networkDevice{
			Device:     fields[0],
			Type:       fields[1],
			State:      fields[2],
			Connection: fields[3],
		}
		devices = append(devices, device)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return devices, nil
}

/*
* 网卡详情:
*   nmcli device show eth0
*
 */
func GetNmcliDeviceShow(c *gin.Context, ruleEngine typex.RuleX) {
	ifaceName, _ := c.GetQuery("iface")
	interfaces, err := ossupport.GetAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ok := false
	for _, iface := range interfaces {
		if iface.Name == ifaceName {
			ok = true
			break
		}
	}
	if !ok {
		c.JSON(common.HTTP_OK, common.Error("interface not exists"))
		return
	}
	cmd := exec.Command("nmcli", "device", "show", "eth0")
	nmcliOutput, err := cmd.Output()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	device, err := parseNmcliDeviceShow(string(nmcliOutput))
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(device))
}

// GENERAL.DEVICE:                 eth0
// GENERAL.TYPE:                   ethernet
// GENERAL.HWADDR:                 02:81:5E:DF:D4:81
// GENERAL.MTU:                    1500
// GENERAL.STATE:                  100 (connected)
// GENERAL.CONNECTION:             eth0
// GENERAL.CON-PATH:               /org/freedesktop/NetworkManager/ActiveConnection/2
// WIRED-PROPERTIES.CARRIER:       on
// IP4.ADDRESS[1]:                 192.168.1.185/24
// IP4.GATEWAY:                    192.168.1.1
// IP4.ROUTE[1]:                   dst = 0.0.0.0/0, nh = 192.168.1.1, mt = 101
// IP4.ROUTE[2]:                   dst = 192.168.1.0/24, nh = 0.0.0.0, mt = 101
// IP4.DNS[1]:                     192.168.1.1
// IP6.ADDRESS[1]:                 fe80::9460:7480:61a9:cbd2/64
// IP6.GATEWAY:                    --
// IP6.ROUTE[1]:                   dst = ff00::/8, nh = ::, mt = 256, table=255
// IP6.ROUTE[2]:                   dst = fe80::/64, nh = ::, mt = 256
// IP6.ROUTE[3]:                   dst = fe80::/64, nh = ::, mt = 101

type networkDeviceDetail struct {
	Device      string `json:"device"`
	Type        string `json:"type"`
	HWAddr      string `json:"hwAddr"`
	MTU         int    `json:"mtu"`
	State       string `json:"state"`
	Connection  string `json:"connection"`
	Carrier     string `json:"carrier"`
	IPv4Addr    string `json:"ipv4Addr"`
	IPv4Gateway string `json:"ipv4Gateway"`
	IPv4DNS     string `json:"ipv4Dns"`
	IPv6Addr    string `json:"ipv6Addr"`
	IPv6Gateway string `json:"ipv6Gateway"`
}

// parseNMCLIOutput 解析 nmcli 输出
// nmcli device show
func parseNmcliDeviceShow(output string) (*networkDeviceDetail, error) {
	lines := strings.Split(output, "\n")

	device := &networkDeviceDetail{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "GENERAL.DEVICE:":
			device.Device = fields[1]
		case "GENERAL.TYPE:":
			device.Type = fields[1]
		case "GENERAL.HWADDR:":
			device.HWAddr = fields[1]
		case "GENERAL.MTU:":
			device.MTU = parseInt(fields[1])
		case "GENERAL.STATE:":
			device.State = fields[1]
		case "GENERAL.CONNECTION:":
			device.Connection = fields[1]
		case "WIRED-PROPERTIES.CARRIER:":
			device.Carrier = fields[1]
		case "IP4.ADDRESS[1]:":
			device.IPv4Addr = fields[1]
		case "IP4.GATEWAY:":
			device.IPv4Gateway = fields[1]
		case "IP4.DNS[1]:":
			device.IPv4DNS = fields[1]
		case "IP6.ADDRESS[1]:":
			device.IPv6Addr = fields[1]
		case "IP6.GATEWAY:":
			device.IPv6Gateway = fields[1]
		}
	}

	return device, nil
}

// parseInt 将字符串转换为整数，如果失败返回 0
func parseInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}
