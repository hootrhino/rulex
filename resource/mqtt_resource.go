package resource

import (
	"fmt"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
type mqttConfig struct {
	Host     string `json:"host" validate:"required" title:"服务地址" info:""`
	Port     int    `json:"port" validate:"required" title:"服务端口" info:""`
	Topic    string `json:"topic" validate:"required" title:"消息来源" info:""`
	ClientId string `json:"clientId" validate:"required" title:"客户端ID" info:""`
	Username string `json:"username" validate:"required" title:"连接账户" info:""`
	Password string `json:"password" validate:"required" title:"连接密码" info:""`
}

//
type MqttInEndResource struct {
	typex.XStatus
	client mqtt.Client
}

func NewMqttInEndResource(inEndId string, e typex.RuleX) typex.XResource {
	m := new(MqttInEndResource)
	m.PointId = inEndId
	m.RuleEngine = e
	return m
}

func (mm *MqttInEndResource) Start() error {
	config := mm.RuleEngine.GetInEnd(mm.PointId).Config
	var mainConfig mqttConfig

	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		mm.RuleEngine.Work(mm.RuleEngine.GetInEnd(mm.PointId), string(msg.Payload()))

	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Mqtt InEnd Connected Success")
		client.Subscribe(mainConfig.Topic, 2, nil)
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Warnf("Connect lost: %v, try to reconnect\n", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mainConfig.Host, mainConfig.Port))
	opts.SetClientID(mainConfig.ClientId)
	opts.SetUsername(mainConfig.Username)
	opts.SetPassword(mainConfig.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	opts.SetMaxReconnectInterval(5 * time.Second)
	mm.client = mqtt.NewClient(opts)
	if token := mm.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}

}

func (mm *MqttInEndResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (m *MqttInEndResource) OnStreamApproached(data string) error {
	return nil
}
func (mm *MqttInEndResource) Stop() {
	mm.client.Disconnect(0)
	mm = nil
}
func (mm *MqttInEndResource) Reload() {

}
func (mm *MqttInEndResource) Pause() {

}
func (mm *MqttInEndResource) Status() typex.ResourceState {
	if mm.client != nil {
		if mm.client.IsConnected() {
			return typex.UP
		} else {
			return typex.DOWN
		}
	} else {
		return typex.DOWN
	}

}

func (mm *MqttInEndResource) Register(inEndId string) error {
	mm.PointId = inEndId
	return nil
}

func (mm *MqttInEndResource) Test(inEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
}

func (mm *MqttInEndResource) Enabled() bool {
	return mm.Enable
}
func (mm *MqttInEndResource) Details() *typex.InEnd {
	return mm.RuleEngine.GetInEnd(mm.PointId)
}
func (*MqttInEndResource) Driver() typex.XExternalDriver {
	return nil
}
func (*MqttInEndResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("MQTT", "", mqttConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

//
// 拓扑
//
func (*MqttInEndResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
