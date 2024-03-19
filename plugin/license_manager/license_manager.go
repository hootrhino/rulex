package licensemanager

/*
*
* 证书管理器
 */

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

type LocalLicense struct {
	DeviceID          string // 设备生产序列号
	AuthorizeAdmin    string // 证书签发人
	AuthorizePassword string // 证书签发人密钥
	BeginAuthorize    int64  // 证书授权开始时间
	EndAuthorize      int64  // 证书授权结束时间
	MAC               string // 设备硬件MAC地址，一般取以太网卡
	License           string // 公钥, 发给用户设备
}

func (d LocalLicense) ValidateTime() bool {
	Now := time.Now().UnixNano()
	V := d.EndAuthorize - Now
	if (d.BeginAuthorize > Now) && (V <= 0) {
		return false
	}
	return true
}

// 00001 & rhino & hoot & FF:FF:FF:FF:FF:FF & 0 & 0
func ParseAuthInfo(info string) (LocalLicense, error) {
	LocalLicense := LocalLicense{}
	ss := strings.Split(info, "&")
	if len(ss) == 6 {
		BeginAuthorize, err1 := strconv.ParseInt(ss[4], 10, 64)
		if err1 != nil {
			return LocalLicense, err1
		}
		EndAuthorize, err2 := strconv.ParseInt(ss[5], 10, 64)
		if err2 != nil {
			return LocalLicense, err2
		}
		LocalLicense.DeviceID = ss[0]
		LocalLicense.AuthorizeAdmin = ss[1]
		LocalLicense.AuthorizePassword = ss[2]
		LocalLicense.MAC = ss[3]
		LocalLicense.BeginAuthorize = BeginAuthorize
		LocalLicense.EndAuthorize = EndAuthorize
		return LocalLicense, nil
	}
	return LocalLicense, fmt.Errorf("failed parse:%s", info)
}

// LicenseManager 证书管理
type LicenseManager struct {
	localLicense LocalLicense
}

/*
*
* 开源版本默认是给一个授权证书
*
 */
func NewLicenseManager(r typex.RuleX) *LicenseManager {
	return &LicenseManager{
		localLicense: LocalLicense{},
	}
}

func (l *LicenseManager) Init(section *ini.Section) error {
	license_path, err1 := section.GetKey("license_path")
	errMsg := "License loading failed. Your License may not be compliant."
	if err1 != nil {
		glogger.GLogger.Fatal()
		os.Exit(0)
	}
	key_path, err := section.GetKey("key_path")
	if err != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	licBytesB64, err := os.ReadFile(license_path.String())
	if err != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	keyBytes, err := os.ReadFile(key_path.String())
	if err != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	licBytes, err := base64.StdEncoding.DecodeString(string(licBytesB64))
	if err != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	adminSalt, err := RSADecrypt(licBytes, keyBytes)
	if err != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	LocalLicense, err := ParseAuthInfo(string(adminSalt))
	if err != nil {
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}

	LocalLicense.License = string(licBytesB64)
	l.localLicense = LocalLicense
	T1 := time.Unix(LocalLicense.BeginAuthorize, 0)
	T2 := time.Unix(LocalLicense.EndAuthorize, 0)
	T1s := T1.Format("2006-01-02 15:04:05")
	T2s := T2.Format("2006-01-02 15:04:05")
	//
	if LocalLicense.BeginAuthorize == 0 && LocalLicense.EndAuthorize == 0 {
		glogger.GLogger.Info("This is Indefinite use version.")
	} else {
		if !LocalLicense.ValidateTime() {
			glogger.GLogger.Fatalf("License has expired, Valid from %s to %s", T1s, T2s)
			os.Exit(0)
		}
		// get local mac .....
	}
	Tip := "Indefinite use version, You can use it till the end of time"
	if LocalLicense.BeginAuthorize == 0 {
		T1s = Tip
	}
	if LocalLicense.EndAuthorize == 0 {
		T2s = Tip
	}
	fmt.Println("[∫∫] Load Local License Success <'v'>")
	fmt.Println("|>>| Device ID:", LocalLicense.DeviceID)
	fmt.Println("|>>| Authorize Admin:", LocalLicense.AuthorizeAdmin)
	fmt.Println("|>>| Authorize Password:", LocalLicense.AuthorizePassword)
	fmt.Println("|>>| Local MAC:", LocalLicense.MAC)
	fmt.Println("|>>| Begin Authorize:", T1s)
	fmt.Println("|>>| End   Authorize:", T2s)
	fmt.Println("|>>| License Content:", LocalLicense.License)
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
		UUID:     "LicenseManager",
		Name:     "LicenseManager",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "HootRhinoTeam",
		Email:    "HootRhinoTeam@hootrhino.com",
		License:  "AGPL",
	}
}
