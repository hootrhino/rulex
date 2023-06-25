package icmpsender

import (
	"fmt"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type ICMPSender struct {
	uuid    string
	pinging bool
}

func NewICMPSender() *ICMPSender {
	return &ICMPSender{
		uuid:    "ICMPSender",
		pinging: false,
	}
}

func (dm *ICMPSender) Init(config *ini.Section) error {
	return nil
}

func (dm *ICMPSender) Start(typex.RuleX) error {
	return nil
}
func (dm *ICMPSender) Stop() error {
	return nil
}

func (hh *ICMPSender) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "ICMPSender",
		Version:  "v1.0.0",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

/*
*
* 服务调用接口
*
 */
func (cs *ICMPSender) Service(arg typex.ServiceArg) typex.ServiceResult {
	// ping 8.8.8.8
	Fields := logrus.Fields{
		"logType":    "ICMPSenderPing",
		"pluginUUID": "ICMPSender",
	}
	if arg.Name == "ping" {
		if cs.pinging {
			glogger.GLogger.WithFields(Fields).Info("ICMPSender pinging now:", arg.Args)
			return typex.ServiceResult{Out: []map[string]interface{}{}}
		}
	}

	switch tt := arg.Args.(type) {
	case []interface{}:
		if len(tt) < 1 {
			break
		}
		cs.pinging = true
		for i := 0; i < 5; i++ {
			switch ip := tt[0].(type) {
			case string:
				if Duration, err := pingQ(ip, 3000*time.Millisecond); err != nil {
					glogger.GLogger.WithFields(Fields).Info(fmt.Sprintf(
						"[Count:%d] Ping Error:%s", i,
						err.Error()))
				} else {
					glogger.GLogger.WithFields(Fields).Info(fmt.Sprintf(
						"[Count:%d] Ping Reply From %s: time=%v ms TTL=128", i,
						tt, Duration))
				}
				time.Sleep(1 * time.Second)
			}

		}
		cs.pinging = false
	default:
		{
			glogger.GLogger.WithFields(Fields).Info("Unknown service name:", arg.Name)
		}
	}

	return typex.ServiceResult{Out: []map[string]interface{}{}}
}
