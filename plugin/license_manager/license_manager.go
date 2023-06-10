package licensemanager

/*
*
* 证书管理器
*
 */

import (
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

type LicenseManager struct {
	uuid string
}

func NewLicenseManager() *LicenseManager {
	return &LicenseManager{
		uuid: "LicenseManager",
	}
}

func (dm *LicenseManager) Init(config *ini.Section) error {
	return nil
}

func (dm *LicenseManager) Start(typex.RuleX) error {
	return nil
}
func (dm *LicenseManager) Stop() error {
	return nil
}

func (hh *LicenseManager) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "LicenseManager",
		Version:  "v0.0.1",
		Homepage: "www.github.com/hootrhino/rulex",
		HelpLink: "www.github.com/hootrhino/rulex",
		Author:   "wwhai",
		Email:    "13594448678@163.com",
		License:  "MIT",
	}
}

/*
 *
 * 服务调用接口
 *
 */
func (cs *LicenseManager) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
