package service

import (
	"fmt"
	"os/exec"

	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

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

/*
*
* 永远只有一个配置,eth0
*
 */
func GetEth0Config() (model.MNetworkConfig, error) {
	MNetworkConfig := model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=?", "eth0").
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}

/*
*
* 永远只有一个配置,eth1
*
 */
func GetEth1Config() (model.MNetworkConfig, error) {
	MNetworkConfig := model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=?", "eth1").
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}

/*
*
* 永远只更新id=0的
*
 */
func UpdateEth0Config(MNetworkConfig model.MNetworkConfig) error {
	Model := model.MNetworkConfig{}
	return interdb.DB().
		Model(Model).
		Where("interface=?", "eth0").
		Updates(MNetworkConfig).Error
}

/*
*
* 永远只更新id=1的
*
 */
func UpdateEth1Config(MNetworkConfig model.MNetworkConfig) error {
	Model := model.MNetworkConfig{}
	return interdb.DB().
		Model(Model).
		Where("interface=?", "eth1").
		Updates(MNetworkConfig).Error
}

/*
*
* 检查一下是否已经初始化过了，避免覆盖配置
*
 */
func CheckIfAlreadyInitNetWorkConfig() bool {
	sql := `SELECT count(*) FROM m_network_configs;`
	count := 0
	err := interdb.DB().Raw(sql).Find(&count).Error
	if err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

/*
*
  - 清空表:DELETE FROM table_name;
    DELETE FROM sqlite_sequence WHERE name='m_network_configs';

*
*/
func TruncateConfig() error {
	sql := `DELETE FROM m_network_configs;DELETE FROM sqlite_sequence WHERE name='m_network_configs';`
	count := 0
	err := interdb.DB().Raw(sql).Find(&count).Error
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

/*
*
* 初始化网卡配置参数
*
 */
func InitNetWorkConfig() error {

	// 默认给DHCP
	dhcp0 := true
	dhcp1 := false
	eth0 := model.MNetworkConfig{
		Interface: "eth0",
		Address:   "192.168.1.100",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.1.1",
		DNS: model.StringList{
			"8.8.8.8",
			"114.114.114.114",
		},
		DHCPEnabled: &dhcp0,
	}
	eth1 := model.MNetworkConfig{
		Interface: "eth1",
		Address:   "192.168.64.100",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.64.1",
		DNS: model.StringList{
			"8.8.8.8",
			"114.114.114.114",
		},
		DHCPEnabled: &dhcp1,
	}
	var err error
	err = interdb.DB().Where("interface=? and id=1", "eth0").FirstOrCreate(&eth0).Error
	if err != nil {
		return err
	}
	err = interdb.DB().Where("interface=? and id=2", "eth1").FirstOrCreate(&eth1).Error
	if err != nil {
		return err
	}
	return nil
}

/*
*
* 应用最新配置
* sudo netplan apply
*
 */
func NetplanApply() error {
	cmd := exec.Command("netplan", "apply")
	cmd.Dir = "/etc/netplan/"
	cmd.Env = append(cmd.Env,
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		glogger.GLogger.Error(err, string(out))
		return err
	}
	return nil
}

// RestartNetworkManager 用于重启 NetworkManager 服务
func RestartNetworkManager() error {
	cmd := exec.Command("systemctl", "restart", "NetworkManager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(err.Error() + ":" + string(output))
	}
	return nil
}

/*
*
* Ubuntu16 ETC 配置应用
* sudo service networking restart
 */
func EtcApply() error {
	cmd := exec.Command("sh", "-c", `service networking restart`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(err.Error() + ":" + string(output))
	}
	return nil
}

/*
*
* 匹配: /etc/network/interfaces
*
 */
func GetAllNetConfig() ([]model.MNetworkConfig, error) {
	// 查出前两个网卡的配置
	ethCfg := []model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=? or interface=?", "eth0", "eth1").
		Find(&ethCfg).Error
	return ethCfg, err
}
