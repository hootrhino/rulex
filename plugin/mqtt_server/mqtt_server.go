package mqttserver

import (
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	mqttServer "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"gopkg.in/ini.v1"
)

const (
	defaultTransport string = "tcp"
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
	uuid       string
}

func NewMqttServer() typex.XPlugin {
	return &MqttServer{
		Host:    "127.0.0.1",
		Port:    1884,
		clients: map[string]*events.Client{},
		uuid:    "RULEX-MqttServer",
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

	if err := server.AddListener(tcpPort, &listeners.Config{
		Auth: &AuthController{},
	}); err != nil {
		return err
	}
	if err := server.Serve(); err != nil {
		return err
	}

	s.mqttServer = server
	s.mqttServer.Events.OnConnect = func(client events.Client, packet events.Packet) {
		s.clients[client.ID] = &client
		glogger.GLogger.Debugf("Client connected:%s", client.ID)
	}
	s.mqttServer.Events.OnDisconnect = func(client events.Client, err error) {
		if s.clients[client.ID] != nil {
			delete(s.clients, client.ID)
			glogger.GLogger.Debugf("Client disconnected:%s", client.ID)
		}
	}
	s.mqttServer.Events.OnMessage = func(c events.Client, p events.Packet) (events.Packet, error) {
		glogger.GLogger.Debug("OnMessage:", c.ID, c.Username, p.TopicName, p.Payload)
		return p, nil
	}
	glogger.GLogger.Infof("MqttServer start at [%s:%v] successfully", s.Host, s.Port)
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
		UUID:     s.uuid,
		Name:     "Light Weight MqttServer",
		Version:  "0.0.1",
		Homepage: "www.github.com/hootrhino/rulex",
		HelpLink: "www.github.com/hootrhino/rulex",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}


/*
*
* 认证器, 目前只做了个简单的认证机制：password=md5(clientid+username)
*
 */
type AuthController struct {
}

func (*AuthController) Authenticate(user, password []byte) bool {
	glogger.GLogger.Debug("Client require Authenticate:", user, string(password))
	return true
}
func (*AuthController) ACL(user []byte, topic string, write bool) bool {
	glogger.GLogger.Debug("Client require ACL:", topic, write)
	return true
}
