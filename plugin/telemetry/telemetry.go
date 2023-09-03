package telemetry

import (
	"encoding/json"
	"net"
	"runtime"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/ini.v1"
)

var _ typex.XPlugin = (*Telemetry)(nil)

type Telemetry struct {
	uuid string
}

type TelemetryConfig struct {
	Enable bool   `ini:"enable"`
	Addr   string `ini:"addr"`
}

// 参数为外部配置
func (*Telemetry) Init(sec *ini.Section) error {
	// 加载配置
	var conf TelemetryConfig
	if err := utils.InIMapToStruct(sec, &conf); err != nil {
		return nil
	}

	if !conf.Enable || len(conf.Addr) == 0 {
		return nil
	}

	// 发起UDP连接
	conn, err := net.Dial("udp", conf.Addr)
	if err != nil {
		// 日志
		glogger.GLogger.Error("plugin.telemetry dail", err.Error())
		// fmt.Println("plugin.telemetry dail", err. common.Error())
		return nil
	}

	// 异步处理，避免阻塞主线程
	go func(conn net.Conn) {
		defer conn.Close()

		// 加载硬件信息
		var info struct {
			OS   string `json:"os"`
			Arch string `json:"arch"`
		}
		info.OS = runtime.GOOS
		info.Arch = runtime.GOARCH

		data, _ := json.Marshal(&info)

		//  发出数据，间隔300ms发出，避免网络阻塞
		for i := 0; i < 5; i++ {
			_, err = conn.Write(data)
			if err != nil {
				glogger.GLogger.Error("plugin.telemetry send", err.Error())
				// fmt.Println("plugin.telemetry send", err. common.Error())
				return
			}
			time.Sleep(300 * time.Millisecond)
		}
	}(conn)

	return nil
}
func (*Telemetry) Start(typex.RuleX) error {
	return nil
}

// 对外提供一些服务
func (*Telemetry) Service(typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
func (*Telemetry) Stop() error {
	return nil
}

func (p *Telemetry) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     "BUSINESS_TELEMETRY",
		Name:     "Business Telemetry",
		Version:  "v0.0.1",
		Homepage: "https://github.com/dropliu/rulex",
		HelpLink: "https://github.com/dropliu/rulex",
		Author:   "dropliu",
		Email:    "13594448678@163.com",
		License:  "MIT",
	}
}
