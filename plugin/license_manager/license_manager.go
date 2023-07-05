package licensemanager

/*
*
* 证书管理器
 */

import (
	"crypto/rsa"
	"crypto/x509/pkix"
	"sync/atomic"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/ini.v1"
)

// LicenseManager 证书管理
type LicenseManager struct {
	uuid string
	cert atomic.Value // 证书信息(*Certificate)
	conf LicenseConfig
}

func NewLicenseManager() *LicenseManager {
	return &LicenseManager{
		uuid: "LicenseManager",
	}
}

// LicenseConfig 证书配置
type LicenseConfig struct {
	LocalAddr  string `ini:"local_addr"`  // 本地存放证书的路径
	RemoteAddr string `ini:"remote_addr"` // 远程拉取证书的路径
	NetName    string `ini:"net_name"`    // 网络名('eth0')
}

// Certificate 证书信息
type Certificate struct {
	Raw       string    `json:"raw"`
	Issuer    pkix.Name `json:"issuer"`  // 授权方
	Subject   pkix.Name `json:"subject"` // 证书授权信息
	NotAfter  time.Time
	NotBefore time.Time
	PublicKey *rsa.PublicKey `json:"key"` // 加密公钥
}

// Init 初始化LicenseManager
func (l *LicenseManager) Init(section *ini.Section) error {
	// 解析配置
	_ = utils.InIMapToStruct(section, &l.conf)

	// 先尝试加载证书
	l.reload(false)

	// 加载失败，设置定时器1小时再检查
	if l.cert.Load() == nil {
		time.AfterFunc(time.Hour, func() { l.reload(true) })
	}

	return nil
}

func (l *LicenseManager) reload(quit bool) *Certificate {
	// 已有证书，直接返回
	cert, _ := l.cert.Load().(*Certificate)
	if cert != nil {
		return cert
	}

	// 尝试加载证书
	cert, err := loadAndVerifyCert(l.conf)
	if err == nil {
		if !l.cert.CompareAndSwap(nil, cert) {
			cert, _ = l.cert.Load().(*Certificate)
		}
		return cert
	}

	// 日志输出
	glogger.GLogger.Error("license_manager:", err.Error())
	// fmt.Println("license_manager:", err. common.Error())

	// 加载失败并退出
	if quit {
		glogger.GLogger.Fatal("need a valid certificate")
		// panic("need a valid certificate")
	}

	return nil
}

// Start 未实现
func (dm *LicenseManager) Start(typex.RuleX) error {
	return nil
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
		Homepage: "https://github.com/dropliu/rulex",
		HelpLink: "https://github.com/dropliu/rulex",
		Author:   "dropliu",
		Email:    "13594448678@163.com",
		License:  "MIT",
	}
}

/*
 *
 * 服务调用接口
 *
 */
func (l *LicenseManager) Service(arg typex.ServiceArg) typex.ServiceResult {
	switch arg.Name {
	case "info":
		var info CertificateInfo
		if cert, _ := l.cert.Load().(*Certificate); cert != nil {
			info = Certificate2Info(cert)
		}
		return typex.ServiceResult{Out: info}
	case "verify":
		var info CertificateInfo
		if cert := l.reload(false); cert != nil {
			info = Certificate2Info(cert)
		}
		return typex.ServiceResult{Out: info}
	default:
	}

	return typex.ServiceResult{}
}

type CertificateInfo struct {
	Authorized bool   `json:"authorized"` // 是否已授权？(当为false时，其它字段为空)
	Content    string `json:"content"`    // 证书原始内容
	StartDate  string `json:"start_date"` // 开始日期
	EndDate    string `json:"end_date"`   // 结束日期
	Valid      bool   `json:"valid"`      // 证书是否有效？
	Issuer     string `json:"issuer"`
	Subject    string `json:"subject"`
}

func Certificate2Info(cert *Certificate) CertificateInfo {
	return CertificateInfo{
		Authorized: true,
		Content:    cert.Raw,
		StartDate:  cert.NotBefore.Format(time.DateOnly),
		EndDate:    cert.NotAfter.Format(time.DateOnly),
		Valid:      time.Since(cert.NotAfter) < 0,
		Issuer:     cert.Issuer.CommonName,
		Subject:    cert.Subject.CommonName,
	}
}
