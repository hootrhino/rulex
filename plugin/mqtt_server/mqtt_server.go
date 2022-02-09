package mqttserver

import (
	"fmt"
	"rulex/typex"

	mqttServer "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/ngaut/log"
)

const (
	defaultPort      int    = 2883
	defaultTransport string = "tcp"
)

type MqttServer struct {
	mqttServer *mqttServer.Server
}

func NewMqttServer() typex.XPlugin {
	return &MqttServer{}
}

func (s *MqttServer) Init(cfg interface{}) error {
	return nil
}

func (s *MqttServer) Start() error {
	server := mqttServer.New()
	tcpPort := listeners.NewTCP(defaultTransport, fmt.Sprintf(":%v", defaultPort))

	if err := server.AddListener(tcpPort, nil); err != nil {
		return err
	}

	if err := server.Serve(); err != nil {
		return err
	}

	s.mqttServer = server
	log.Info("MqttServer start at [0.0.0.0:2883] successfully")
	return nil
}

func (s *MqttServer) Stop() error {
	if s.mqttServer != nil {
		return s.mqttServer.Close()
	} else {
		return nil
	}

}

func (s *MqttServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "Light Weight MqttServer",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
