package mqttserver

import (
	"context"
	"net"
	"rulex/typex"

	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/ngaut/log"
	"go.uber.org/zap"
)

type MqttServer struct {
	mqttServer server.Server
}

func NewMqttServer() typex.XPlugin {
	return &MqttServer{}
}

func (s *MqttServer) Init() error {
	return nil
}

func (s *MqttServer) Start() error {
	tcpPort, err := net.Listen("tcp", ":1883")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	logger, _ := zap.NewDevelopment()
	mqttServer := server.New(
		server.WithTCPListener(tcpPort),
		server.WithLogger(logger),
		server.WithConfig(config.DefaultConfig()),
	)
	if err := mqttServer.Run(); err != nil {
		return err
	}
	s.mqttServer = mqttServer
	log.Info("MqttServer start successfully")
	return nil
}

func (s *MqttServer) Stop() error {
	log.Info("MqttServer stop successfully")
	return s.mqttServer.Stop(context.Background())
}

func (s *MqttServer) XPluginMetaInfo() typex.XPluginMetaInfo {
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
