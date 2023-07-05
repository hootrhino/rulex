package licensemanager

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hootrhino/rulex/glogger"
)

const defaultCertPath = "./rulex.pem"

func loadAndVerifyCert(conf LicenseConfig) (*Certificate, error) {
	path := conf.LocalAddr
	if len(path) == 0 {
		path = defaultCertPath
	}
	// 加载证书
	cert, isLocal, err := loadCert(path, conf.RemoteAddr, conf.NetName)
	if err != nil {
		return nil, err
	}

	if !isLocal {
		// 第一次拉到证书，联网校验
		if err := onlineVerifyCert(conf.RemoteAddr, conf.NetName, cert.PublicKey); err != nil {
			return nil, err
		}
		// 校验成功，证书写入到本地
		writeToFile(path, cert.Raw)
	} else {
		// 本地已有证书，离线校验即可
		if err := offlineVerifyCert(cert, conf.NetName); err != nil {
			// 验证失败，移除证书
			os.Remove(path)
			return nil, err
		}
	}
	return cert, nil
}

func writeToFile(path string, s string) {
	file, err := os.Create(path)
	if err != nil {
		glogger.GLogger.Error("license_manager:writeToFile", err.Error())
		// fmt.Println("license_manager:writeToFile", err. common.Error())
		return
	}
	defer file.Close()
	file.Write(s2b(s))
	file.Sync()
}

func loadCert(local, remote, netName string) (*Certificate, bool, error) {
	cert, err := loadLocalCert(local)
	if err == nil {
		return cert, true, nil
	}

	cert, err1 := pullRemoteCert(remote, netName)
	if err1 == nil {
		// 写入到本地
		return cert, false, nil
	}

	return nil, false, fmt.Errorf("load cert from local '%s': %s, load cert from remote '%s': %s",
		local, err.Error(), remote, err1.Error())
}

func loadLocalCert(path string) (*Certificate, error) {
	// 没有配置路径，使用默认路径
	if len(path) == 0 {
		path = defaultCertPath
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读证书文件
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cert, err := decodeCert(data)

	return cert, err
}

// DeviceInfo 硬件信息
type DeviceInfo struct {
	MAC  string `json:"mac"` // 硬件mac地址
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type CertificateResponse struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// pullRemoteCert 请求远程主机生成证书并下发
/*
POST http://remote.site/xxx/cert
request:
{
	"mac": "xxx",
	"os": "linux",
	"arch": "arm"
}
response:
{
	"ret":0,
	"msg": "",
	"data": "-----BEGIN CERTIFICATE----"
}
*/
func pullRemoteCert(url, netName string) (*Certificate, error) {
	if len(url) == 0 {
		return nil, errors.New("empty url")
	}
	hw := deviceInfo(netName)
	if len(hw.MAC) == 0 {
		return nil, ErrUnknownHardwareAddr
	}

	// 构建请求
	infoBytes, _ := json.Marshal(&hw)
	body := bytes.NewReader(infoBytes)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	// 请求证书
	var cert *Certificate
	err = retryRequest(req, func(resp *http.Response) error {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var cr CertificateResponse
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}

		if cr.Ret != 0 {
			return fmt.Errorf("ret: %d, msg: %s", cr.Ret, cr.Msg)
		}

		// 解析证书
		cert, err = decodeCert(s2b(cr.Data))
		return err
	})

	return cert, err
}

type VerifyCertRequest struct {
	MAC  string `json:"mac"`
	Data string `json:"data"`
}

type VerifyCertResponse struct {
	Ret int    `json:"ret"` // 返回值
	Msg string `json:"msg"` //
}

/*
GET http://remote.site/xxx/cert
request:

	{
		"mac": "xxx", // 硬件地址，标识机器
		"data" "xxx" // 加密后的数据
	}

response:
*/

var (
	ErrUnsupportedPublicKey  = errors.New("unsupported public key type")
	ErrUnknownHardwareAddr   = errors.New("unknown hardware addr")
	ErrMismatchedCertificate = errors.New("mismatched certificate")
)

func onlineVerifyCert(url, netName string, key *rsa.PublicKey) error {
	hw := deviceInfo(netName)
	if len(hw.MAC) == 0 {
		return ErrUnknownHardwareAddr
	}

	// 加密数据
	src := []byte(hw.MAC)
	dst, err := rsa.EncryptPKCS1v15(rand.Reader, key, src)
	if err != nil {
		return err
	}

	// 构建请求
	rc := VerifyCertRequest{
		MAC:  string(src),
		Data: base64.StdEncoding.EncodeToString(dst),
	}

	bts, _ := json.Marshal(&rc)
	body := bytes.NewReader(bts)
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return err
	}

	// 请求验证
	err = retryRequest(req, func(resp *http.Response) error {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var vc VerifyCertResponse
		if err := json.Unmarshal(data, &vc); err != nil {
			return err
		}

		// 检查响应
		if vc.Ret != 0 {
			return fmt.Errorf("ret: %d, msg: %s", vc.Ret, vc.Msg)
		}
		return nil
	})

	return err
}

// retryRequest 重试请求，第1次超时1.5s，每次重试增加0.75s
func retryRequest(req *http.Request, fn func(*http.Response) error) error {
	var err error
	for i := 0; i < 3; i++ {
		client := http.Client{
			Timeout: time.Duration(1500+i*750) * time.Millisecond,
		}
		// 发起请求
		var resp *http.Response
		resp, err = client.Do(req)
		if err != nil {
			continue
		}
		// 处理响应
		if err = fn(resp); err != nil {
			continue
		}

		return nil
	}

	return err
}

func offlineVerifyCert(cert *Certificate, netName string) error {
	var mac string
	inte, err := net.InterfaceByName(netName)
	if err == nil && inte.HardwareAddr != nil {
		mac = inte.HardwareAddr.String()
		md5Byts := md5.Sum(s2b(mac))
		mac = hex.EncodeToString(md5Byts[:])
	}
	if len(mac) == 0 {
		return fmt.Errorf("unknown netName: %s", netName)
	}

	if mac != cert.Subject.CommonName {
		return ErrMismatchedCertificate
	}
	return nil
}

func deviceInfo(name string) DeviceInfo {
	info := DeviceInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	inte, err := net.InterfaceByName(name)
	if err == nil && inte.HardwareAddr != nil {
		info.MAC = inte.HardwareAddr.String()
	}

	return info
}
