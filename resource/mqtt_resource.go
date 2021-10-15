package resource

import (
	"fmt"
	"rulex/typex"
	"time"

	"github.com/ngaut/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
const DEFAULT_CLIENT_ID string = "X_IN_END_CLIENT"
const DEFAULT_USERNAME string = "X_IN_END"
const DEFAULT_PASSWORD string = "X_IN_END"
const DEFAULT_TOPIC string = "$X_IN_END"

//
type MqttInEndResource struct {
	typex.XStatus
	client mqtt.Client
}

func NewMqttInEndResource(inEndId string, e typex.RuleX) *MqttInEndResource {
	m := new(MqttInEndResource)
	m.PointId = inEndId
	m.RuleEngine = e
	return m
}

func (mm *MqttInEndResource) Start() error {

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		if mm.Enable {
			log.Debug("Message payload:", string(msg.Payload()))
			mm.RuleEngine.Work(mm.RuleEngine.GetInEnd(mm.PointId), string(msg.Payload()))
		}
	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Mqtt InEnd Connected Success")
		// TODO support multipul topics
		client.Subscribe(DEFAULT_TOPIC, 2, nil)
		mm.RuleEngine.GetInEnd(mm.PointId).SetState(typex.UP)
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Infof("Connect lost: %v\n", err)
		time.Sleep(5 * time.Second)
		mm.RuleEngine.GetInEnd(mm.PointId).SetState(typex.DOWN)
	}
	config := mm.RuleEngine.GetInEnd(mm.PointId).Config
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", (*config)["server"], (*config)["port"]))
	if (*config)["clientId"] != nil {
		opts.SetClientID((*config)["clientId"].(string))
	} else {
		opts.SetClientID(DEFAULT_CLIENT_ID)
	}
	if (*config)["username"] != nil {
		opts.SetUsername((*config)["username"].(string))
	} else {
		opts.SetUsername(DEFAULT_USERNAME)
	}
	if (*config)["password"] != nil {
		opts.SetPassword((*config)["password"].(string))
	} else {
		opts.SetPassword(DEFAULT_PASSWORD)
	}
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		log.Warn("Client disconnected, Try to reconnect...")
	}
	opts.SetMaxReconnectInterval(5 * time.Second)
	mm.client = mqtt.NewClient(opts)
	mm.Enable = true
	if token := mm.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}

}

func (mm *MqttInEndResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
}

func (mm *MqttInEndResource) Stop() {
	mm.client.Disconnect(0)

}
func (mm *MqttInEndResource) Reload() {

}
func (mm *MqttInEndResource) Pause() {

}
func (mm *MqttInEndResource) Status() typex.ResourceState {
	if mm.client.IsConnected() {
		return typex.UP
	} else {
		return typex.DOWN
	}
}

func (mm *MqttInEndResource) Register(inEndId string) error {
	mm.PointId = inEndId
	return nil
}

func (mm *MqttInEndResource) Test(inEndId string) bool {
	return mm.client.IsConnected()
}

func (mm *MqttInEndResource) Enabled() bool {
	return mm.Enable
}
func (mm *MqttInEndResource) Details() *typex.InEnd {
	return mm.RuleEngine.GetInEnd(mm.PointId)
}
