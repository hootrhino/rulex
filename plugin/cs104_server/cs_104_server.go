package plugin

import (
	"fmt"

	"gopkg.in/ini.v1"

	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

/*
*
* 配置信息,从ini文件里面读取出来的
*
 */
type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
}
type cs104Server struct {
	server  *cs104.Server
	Host    string
	Port    int
	LogMode bool
	uuid    string
}

func NewCs104Server() typex.XPlugin {
	return &cs104Server{uuid: "CS104-SERVER"}
}
func (sf *cs104Server) InterrogationHandler(c asdu.Connect,
	asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	return nil
}
func (sf *cs104Server) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error {
	return nil
}
func (sf *cs104Server) ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error {
	return nil
}
func (sf *cs104Server) ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error {
	return nil
}
func (sf *cs104Server) ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error {
	return nil
}
func (sf *cs104Server) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error {
	return nil
}
func (sf *cs104Server) ASDUHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

// ---------------------------------------------------------------------------
//
// ---------------------------------------------------------------------------
func (cs *cs104Server) Init(config *ini.Section) error {
	cs.server = cs104.NewServer(&cs104Server{})
	var mainConfig _serverConfig
	if err := utils.InIMapToStruct(config, &mainConfig); err != nil {
		return err
	}
	cs.Host = mainConfig.Host
	cs.Port = mainConfig.Port
	return nil
}

func (cs *cs104Server) Start(typex.RuleX) error {
	cs.server.SetOnConnectionHandler(func(c asdu.Connect) {
		glogger.GLogger.Warn("Connected: ", c.Params())
	})
	cs.server.SetConnectionLostHandler(func(c asdu.Connect) {
		glogger.GLogger.Warn("Disconnected: ", c.Params())
	})
	cs.server.LogMode(cs.LogMode)
	cs.server.ListenAndServer(fmt.Sprintf("%s:%d", cs.Host, cs.Port))
	return nil
}

func (cs *cs104Server) Stop() error {
	return nil
}

func (cs *cs104Server) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     cs.uuid,
		Name:     "IEC104 server Plugin",
		Version:  "v0.0.1",
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
func (cs *cs104Server) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
