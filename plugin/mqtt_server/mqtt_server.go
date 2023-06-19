package mqttserver

import (
	"fmt"
	"sync"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/rs/zerolog"
	"gopkg.in/ini.v1"
)

type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
}
type _topic struct {
	Topic string
}
type MqttServer struct {
	Enable     bool
	Host       string
	Port       int
	mqttServer *mqtt.Server
	clients    map[string]*mqtt.Client
	topics     map[string][]_topic // Topic 订阅表
	ruleEngine typex.RuleX
	uuid       string
}

func NewMqttServer() typex.XPlugin {
	return &MqttServer{
		Host:    "127.0.0.1",
		Port:    1884,
		clients: map[string]*mqtt.Client{},
		topics:  map[string][]_topic{},
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
	zlog := zerolog.New(glogger.GLogger.Writer()).With().Logger()

	server := mqtt.New(&mqtt.Options{Logger: &zlog})
	tcp := listeners.NewTCP("node1", fmt.Sprintf("%v:%v", s.Host, s.Port), nil)
	if err := server.AddListener(tcp); err != nil {
		return err
	}
	if err := server.Serve(); err != nil {
		return err
	}
	//
	// 本地服务器
	//
	s.mqttServer = server
	server.AddHook(&ahooks{s: s}, nil)
	server.AddHook(&mhooks{s: s, locker: sync.Mutex{}}, nil)
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
		Version:  "v2.0.0",
		Homepage: "https://github.com/lion-brave",
		HelpLink: "https://github.com/lion-brave",
		Author:   "liyong",
		Email:    "liyong@gmail.com",
		License:  "MIT",
	}
}

/*
*
* 认证器
*
 */

type mhooks struct {
	mqtt.HookBase
	s      *MqttServer
	locker sync.Mutex
}

func (h *mhooks) ID() string {
	return "events-hooks"
}

func (h *mhooks) Provides(b byte) bool {
	return true
}

func (h *mhooks) OnConnect(client *mqtt.Client, pk packets.Packet) {
	h.locker.Lock()
	h.s.clients[client.ID] = client
	h.locker.Unlock()
	glogger.GLogger.Debugf("client OnConnect:[%v] %v", client.ID, string(client.Properties.Username))

}

func (h *mhooks) OnDisconnect(client *mqtt.Client, err error, expire bool) {
	if h.s.clients[client.ID] != nil {
		h.locker.Lock()
		delete(h.s.clients, client.ID)
		delete(h.s.topics, client.ID)
		h.locker.Unlock()
		glogger.GLogger.Debugf("Client disconnected:%s", client.ID)
	}
}

func (h *mhooks) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	glogger.GLogger.Debugf("client OnPublish:[%v]=%v", pk.TopicName, string(pk.Payload))
	return pk, nil
}

// ahooks is an authentication hook which allows connection access
// for all users and read and write access to all topics.
type ahooks struct {
	mqtt.HookBase
	s *MqttServer
}

// ID returns the ID of the hook.
func (h *ahooks) ID() string {
	return "auth-hooks"
}

// Provides indicates which hook methods this hook provides.
func (h *ahooks) Provides(b byte) bool {
	return true
}

// OnConnectAuthenticate returns true/allowed for all requests.
func (h *ahooks) OnConnectAuthenticate(client *mqtt.Client, pk packets.Packet) bool {
	glogger.GLogger.Debugf("OnAuthenticate:[%v],[%v],[%v]",
		client.ID, string(client.Properties.Username), string(pk.Connect.Password))
	return true
}

// OnACLCheck returns true/allowed for all checks.
func (h *ahooks) OnACLCheck(client *mqtt.Client, topic string, write bool) bool {
	glogger.GLogger.Debugf("OnACLCheck:[%v],[%v],[%v]",
		client.ID, string(client.Properties.Username), topic)
	_, ok := h.s.topics[client.ID]
	if !ok {
		h.s.topics[client.ID] = []_topic{{Topic: topic}}
	} else {
		h.s.topics[client.ID] = append(h.s.topics[client.ID], _topic{Topic: topic})
	}
	return true
}
