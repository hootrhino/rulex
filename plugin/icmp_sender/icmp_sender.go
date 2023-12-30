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
		Name:     "ICMP Sender",
		Version:  "v1.0.0",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "HootRhinoTeam",
		Email:    "HootRhinoTeam@hootrhino.com",
		License:  "AGPL",
	}
}

/*
*
* 服务调用接口
*
 */
func (icmp *ICMPSender) Service(arg typex.ServiceArg) typex.ServiceResult {
	// ping 8.8.8.8
	Fields := logrus.Fields{
		"topic": "plugin/ICMPSenderPing/ICMPSender",
	}
	out := typex.ServiceResult{Out: []map[string]interface{}{}}
	if icmp.pinging {
		glogger.GLogger.WithFields(Fields).Info("ICMPSender pinging now:", arg.Args)
		return out
	}
	if arg.Name == "ping" {
		icmp.pinging = true
		go func(cs *ICMPSender) {
			defer func() {
				cs.pinging = false
			}()
			select {
			case <-typex.GCTX.Done():
				{
					return
				}
			default:
				{
				}
			}
			switch tt := arg.Args.(type) {
			case []interface{}:
				if len(tt) < 1 {
					break
				}
				for i := 0; i < 5; i++ {
					switch ip := tt[0].(type) {
					case string:
						if Duration, err := pingQ(ip, 1000*time.Millisecond); err != nil {
							glogger.GLogger.WithFields(Fields).Info(fmt.Sprintf(
								"[Count:%d] Ping Error:%s", i, err.Error()))
						} else {
							glogger.GLogger.WithFields(Fields).Info(fmt.Sprintf(
								"[Count:%d] Ping Reply From %s: time=%v ms TTL=128", i, tt, Duration))
						}
						// 300毫秒
						time.Sleep(100 * time.Millisecond)
					}

				}
			default:
				{
					glogger.GLogger.WithFields(Fields).Info("Unknown service name:", arg.Name)
				}
			}
		}(icmp)

	}
	return out
}
