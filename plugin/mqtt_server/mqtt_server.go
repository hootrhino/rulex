package mqttserver

import (
	"fmt"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	mqttServer "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"gopkg.in/ini.v1"
)

const (
	defaultTransport string = "tcp"
	banner           string = `
    __  _______  ____________   _____ __________ _    ____________
   /  |/  / __ \/_  __/_  __/  / ___// ____/ __ \ |  / / ____/ __ \
  / /|_/ / / / / / /   / /_____\__ \/ __/ / /_/ / | / / __/ / /_/ /
 / /  / / /_/ / / /   / /_____/__/ / /___/ _, _/| |/ / /___/ _, _/
/_/  /_/\___\_\/_/   /_/     /____/_____/_/ |_| |___/_____/_/ |_|`
)

type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
}
type MqttServer struct {
	Enable     bool
	Host       string
	Port       int
	mqttServer *mqttServer.Server
	clients    map[string]*events.Client
	ruleEngine typex.RuleX
}

func NewMqttServer() typex.XPlugin {
	return &MqttServer{
		clients: map[string]*events.Client{},
	}
}

func (s *MqttServer) Init(config *ini.Section) error {
	var mainConfig _serverConfig
	if err := utils.InIMapToStruct(config, &mainConfig); err != nil {
		return err
	}
	s.Host = mainConfig.Host
	s.Port = mainConfig.Port
	return nil
}

func (s *MqttServer) Start(r typex.RuleX) error {
	s.ruleEngine = r
	server := mqttServer.New()
	tcpPort := listeners.NewTCP(defaultTransport, fmt.Sprintf(":%v", s.Port))

	if err := server.AddListener(tcpPort, &listeners.Config{}); err != nil {
		return err
	}
	if err := server.Serve(); err != nil {
		return err
	}

	s.mqttServer = server
	s.mqttServer.Events.OnConnect = func(client events.Client, packet events.Packet) {
		s.clients[client.ID] = &client
		glogger.GLogger.Infof("Client connected:%s", client.ID)
	}
	s.mqttServer.Events.OnDisconnect = func(client events.Client, err error) {
		if s.clients[client.ID] != nil {
			delete(s.clients, client.ID)
			glogger.GLogger.Warnf("Client disconnected:%s", client.ID)
		}
	}
	s.mqttServer.Events.OnMessage = func(c events.Client, p events.Packet) (events.Packet, error) {

		return p, nil
	}
	fmt.Println(banner)
	glogger.GLogger.Infof("MqttServer start at [0.0.0.0:%v] successfully", s.Port)
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
		Homepage: "www.github.com/i4de/rulex",
		HelpLink: "www.github.com/i4de/rulex",
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
func (cs *MqttServer) Service(arg typex.ServiceArg) error {
	return nil
}
