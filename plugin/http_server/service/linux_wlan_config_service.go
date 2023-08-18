package service

import (
	"github.com/hootrhino/rulex/glogger"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

/*
*
* 配置WIFI Wlan0
*
 */
func UpdateWlan0Config(MNetworkConfig model.MWifiConfig) error {
	Model := model.MWifiConfig{Interface: "wlan0"}
	return sqlitedao.Sqlite.DB().
		Model(Model).
		Where("interface=? and id = 1", "wlan0").
		Updates(MNetworkConfig).Error
}

/*
*
* 获取Wlan0的配置信息
*
 */
func GetWlan0Config() (model.MWifiConfig, error) {
	MWifiConfig := model.MWifiConfig{}
	err := sqlitedao.Sqlite.DB().
		Find(&MWifiConfig).
		Where("interface=? and id = 1", "wlan0").Error
	return MWifiConfig, err
}

/*
*
* 检查是否设置了WIFI网络
*
 */
func CheckIfAlreadyInitWlanConfig() bool {
	sql := `SELECT count(*) FROM m_wifi_configs;`
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
* 清空WIFI配置表
*
 */
func TruncateWifiConfig() error {
	sql := `DELETE FROM m_wifi_configs;DELETE FROM sqlite_sequence WHERE name='m_wifi_configs';`
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
func InitWlanConfig() error {

	// 默认给DHCP
	wlan0 := model.MWifiConfig{
		Interface: "wlan0",
		SSID:      "example.net",
		Password:  "123456",
		Security:  "wpa2-psk",
	}
	err := sqlitedao.Sqlite.DB().
		Where("interface=? and id=1", "wlan0").
		FirstOrCreate(&wlan0).Error
	if err != nil {
		return err
	}
	return nil
}
