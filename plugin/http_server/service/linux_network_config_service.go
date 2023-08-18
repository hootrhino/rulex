package service

import (
	"os/exec"

	"github.com/hootrhino/rulex/glogger"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
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
	err := sqlitedao.Sqlite.DB().
		Find(&MNetworkConfig).
		Where("interface=?", "eth0").Error
	return MNetworkConfig, err
}

/*
*
* 永远只有一个配置,eth1
*
 */
func GetEth1Config() (model.MNetworkConfig, error) {
	MNetworkConfig := model.MNetworkConfig{Interface: "eth0"}
	err := sqlitedao.Sqlite.DB().
		Find(&MNetworkConfig).
		Where("interface=?", "eth1").Error
	return MNetworkConfig, err
}

/*
*
* 永远只更新id=0的
*
 */
func UpdateEth0Config(MNetworkConfig model.MNetworkConfig) error {
	Model := model.MNetworkConfig{Interface: "eth1"}
	return sqlitedao.Sqlite.DB().
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
	return sqlitedao.Sqlite.DB().
		Model(Model).
		Where("interface=?", "eth1").
		Updates(MNetworkConfig).Error
}

// /*
// *
// * 初始化数据，开局的时候插入一条记录
// *
//  */
// const ubuntu16CfgDHCP string = `
// auto lo
// iface lo inet loopback
// auto eth0
// iface eth0 inet dhcp
// `
// const ubuntu16CfgStatic string = `
// # local loopback
// auto lo
// iface lo inet loopback
// # eth0
// auto eth0
// iface eth0 inet static
//     address 192.168.128.100
//     netmask 255.255.255.0
//     gateway 192.168.128.1
//     dns-nameservers 8.8.8.8
// # eth1
// auto eth1
// iface eth1 inet static
//     address 192.168.12800.100
//     netmask 255.255.255.0
//     gateway 192.168.12800.1
//     dns-nameservers 8.8.8.8 114.114.114.114
// `
// const ubuntuNetplanDHCPCfg string = `
// network:
//   version: 2
//   renderer: networkd
//   ethernets:
//     eth0:
//       dhcp4: true
//     eth1:
//       dhcp4: true
// `
// const ubuntuNetplanStaticCfg string = `
// network:
//   version: 2
//   renderer: networkd
//   ethernets:
//     eth0:
//       addresses:
//         - 192.168.128.100/24
//       gateway4: 192.168.128.1
//       nameservers:
//         addresses: [8.8.8.8, 8.8.4.4]
//     eth1:
//       addresses:
//         - 10.0.0.100/24
//       gateway4: 10.0.0.1
//       nameservers:
//         addresses: [8.8.8.8, 8.8.4.4]
// `

/*
*
* 检查一下是否已经初始化过了，避免覆盖配置
*
 */
func CheckIfAlreadyInitNetWorkConfig() bool {
	sql := `SELECT count(*) FROM m_network_configs;`
	count := 0
	err := sqlitedao.Sqlite.DB().Raw(sql).Find(&count).Error
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
	err := sqlitedao.Sqlite.DB().Raw(sql).Find(&count).Error
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
	eth0 := model.MNetworkConfig{
		Interface: "eth0",
		Address:   "192.168.128.100",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.128.1",
		DNS: model.StringList{
			"8.8.8.8",
			"114.114.114.114",
		},
		DHCPEnabled: false,
	}
	eth1 := model.MNetworkConfig{
		Interface: "eth1",
		Address:   "192.168.128.100",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.128.1",
		DNS: model.StringList{
			"8.8.8.8",
			"114.114.114.114",
		},
		DHCPEnabled: false,
	}
	var err error
	err = sqlitedao.Sqlite.DB().Where("interface=? and id=1", "eth0").FirstOrCreate(&eth0).Error
	if err != nil {
		return err
	}
	err = sqlitedao.Sqlite.DB().Where("interface=? and id=2", "eth1").FirstOrCreate(&eth1).Error
	if err != nil {
		return err
	}
	return nil
}

/*
*
* 应用最新配置
*
 */
func NetplanApply() error {
	cmd := exec.Command("netplan", "apply")
	cmd.Dir = "/etc/netplan/"
	cmd.Env = append(cmd.Env,
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	glogger.GLogger.Info(out)
	return nil
}
