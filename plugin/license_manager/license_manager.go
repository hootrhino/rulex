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

/*
*
* 开源版本默认是给一个授权证书
*
 */
func NewLicenseManager(r typex.RuleX) *LicenseManager {
	return &LicenseManager{
		uuid:          "LicenseManager",
		AdminUsername: "rhino",
		AdminPassword: "hoot",
		CpuId:         "o:p:e:n:r:h:i:n:o",
		EthMac:        "o:p:e:n:h:o:o:t",
	}
}

func (l *LicenseManager) Init(section *ini.Section) error {
	license_path, err1 := section.GetKey("license_path")
	errMsg := "Load License Public Cipher Failed, May be Your License IS Invalid."
	if err1 != nil {
		glogger.GLogger.Fatal()
		os.Exit(0)
	}
	license_key, err2 := section.GetKey("key_path")
	if err2 != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	lic, err1 := os.ReadFile(license_path.String())
	if err1 != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	license := string(lic)
	key, err2 := os.ReadFile(license_key.String())
	if err2 != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	private := string(key)
	localMacSum := SumMd5(fmt.Sprintf("%s,%s", l.CpuId, l.EthMac))
	localAdminSum := SumMd5(fmt.Sprintf("%s,%s", l.AdminUsername, l.AdminPassword))
	adminSalt, err3 := DecryptAES(private, license)
	if err3 != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	if localAdminSum != private {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	if adminSalt != localMacSum {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	glogger.GLogger.Info("Load License Success:", license)
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
