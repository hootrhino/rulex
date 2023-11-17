package licensemanager

/*
*
* 证书管理器
 */

import (
	"fmt"
	"os"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

// LicenseManager 证书管理
type LicenseManager struct {
	uuid          string
	CpuId         string
	EthMac        string
	AdminUsername string
	AdminPassword string
}

//	admin := {
//		Username: "open-rhino",
//		Password: "open-hoot",
//	}
//
// // 出厂的第一台设备
//
//	device := {
//		CpuId:  "o:p:e:n:r:h:i:n:o",
//		EthMac: "o:p:e:n:h:o:o:t",
//	}
func NewLicenseManager(r typex.RuleX) *LicenseManager {
	return &LicenseManager{
		uuid:          "LicenseManager",
		AdminUsername: "rhino",
		AdminPassword: "hoot",
		CpuId:         "o:p:e:n:r:h:i:n:o",
		EthMac:        "o:p:e:n:h:o:o:t",
	}
}

// Init 初始化LicenseManager
func (l *LicenseManager) Init(section *ini.Section) error {
	licence_path, err1 := section.GetKey("licence_path")
	if err1 != nil {
		glogger.GLogger.Fatal("Load License Public Cipher Failed, May be Your License IS Invalid.")
	}
	licence_key, err2 := section.GetKey("key_path")
	if err2 != nil {
		glogger.GLogger.Fatal("Load License Public Cipher Failed, May be Your License IS Invalid.")
	}
	lic, err1 := os.ReadFile(licence_path.String())
	if err1 != nil {
		glogger.GLogger.Fatal("Load License Public Cipher Failed, May be Your License IS Invalid.")
	}
	licence := string(lic)
	glogger.GLogger.Info("licence public key:", licence)
	key, err2 := os.ReadFile(licence_key.String())
	if err2 != nil {
		glogger.GLogger.Fatal("Load License Public Cipher Failed, May be Your License IS Invalid.")
	}
	private := string(key)
	localMacSum := SumMd5(fmt.Sprintf("%s,%s", l.CpuId, l.EthMac))
	localAdminSum := SumMd5(fmt.Sprintf("%s,%s", l.AdminUsername, l.AdminPassword))
	adminSalt, err3 := DecryptAES(private, licence)
	if err3 != nil {
		glogger.GLogger.Fatal("Load License Public Cipher Failed, May be Your License IS Invalid.")

	}
	if localAdminSum != private {
		glogger.GLogger.Fatal("Load License Public Cipher Failed, May be Your License IS Invalid.")
	}
	if adminSalt != localMacSum {
		glogger.GLogger.Fatal("Load License Failed, May be Your License IS Invalid.")
	}
	glogger.GLogger.Info("Load License Success:", licence)
	return nil
}

// Start 未实现
func (dm *LicenseManager) Start(typex.RuleX) error {
	return nil
}
func (dm *LicenseManager) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

// Stop 未实现
func (dm *LicenseManager) Stop() error {
	return nil
}

func (hh *LicenseManager) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "LicenseManager",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "HootRhinoTeam",
		Email:    "HootRhinoTeam@hootrhino.com",
		License:  "AGPL",
	}
}
