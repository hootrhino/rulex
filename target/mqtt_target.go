package target

import (
	"fmt"
	"rulex/typex"
	"time"

	"github.com/ngaut/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
const DEFAULT_CLIENT_ID string = "X_OUT_END_CLIENT_00000000"
const DEFAULT_USERNAME string = "X_OUT_END_CLIENT_00000000"
const DEFAULT_PASSWORD string = "X_OUT_END"
const DEFAULT_SUB_TOPIC string = "X_OUT_END_CLIENT_00000000/SUB"
const DEFAULT_PUB_TOPIC string = "X_OUT_END_CLIENT_00000000/PUB"

//
type MqttOutEndTarget struct {
	typex.XStatus
	client mqtt.Client
}

func NewMqttTarget(e typex.RuleX) typex.XTarget {
	m := new(MqttOutEndTarget)
	m.RuleEngine = e
	return m
}

func (mm *MqttOutEndTarget) Start() error {

	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Mqtt InEnd Connected Success")
		client.Subscribe(DEFAULT_SUB_TOPIC, 2, nil)
		mm.RuleEngine.GetOutEnd(mm.PointId).SetState(typex.UP)
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Errorf("Connect lost: %v\n", err)
		time.Sleep(5 * time.Second)
		mm.RuleEngine.GetOutEnd(mm.PointId).SetState(typex.DOWN)
	}
	config := mm.RuleEngine.GetOutEnd(mm.PointId).Config
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
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		log.Error("Client disconnected, Try to reconnect...")
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
func (m *MqttOutEndTarget) OnStreamApproached(data string) error {
	return nil
}
func (mm *MqttOutEndTarget) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
}

func (mm *MqttOutEndTarget) Stop() {
	mm.client.Disconnect(0)

}
func (mm *MqttOutEndTarget) Reload() {

}
func (mm *MqttOutEndTarget) Pause() {

}
func (mm *MqttOutEndTarget) Status() typex.ResourceState {
	if mm.client.IsConnected() {
		return typex.UP
	} else {
		return typex.DOWN
	}
}

func (mm *MqttOutEndTarget) Register(inEndId string) error {
	mm.PointId = inEndId
	return nil
}

func (mm *MqttOutEndTarget) Test(inEndId string) bool {
	return mm.client.IsConnected()
}

func (mm *MqttOutEndTarget) Enabled() bool {
	return mm.Enable
}
func (mm *MqttOutEndTarget) Details() *typex.OutEnd {
	return mm.RuleEngine.GetOutEnd(mm.PointId)
}

//
//
//
func (m *MqttOutEndTarget) To(data interface{}) error {
	return m.client.Publish(DEFAULT_PUB_TOPIC, 2, false, data).Error()
}
