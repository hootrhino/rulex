package target

import (
	"errors"
	"fmt"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
type mqttConfig struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Topic    string `json:"topic" validate:"required"`
	ClientId string `json:"clientId" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

//
type MqttOutEndTarget struct {
	typex.XStatus
	client mqtt.Client
	Topic  string
}

func NewMqttTarget(e typex.RuleX) typex.XTarget {
	m := new(MqttOutEndTarget)
	m.RuleEngine = e
	return m
}
func (*MqttOutEndTarget) Driver() typex.XExternalDriver {
	return nil
}
func (mm *MqttOutEndTarget) Start() error {
	config := mm.RuleEngine.GetOutEnd(mm.PointId).Config
	var mainConfig mqttConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	opts := mqtt.NewClientOptions()

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		mm.RuleEngine.GetOutEnd(mm.PointId).Target.To(string(msg.Payload()))
	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Mqtt OutEnd Connected Success")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Warnf("Connect lost: %v, try to reconnect\n", err)
		time.Sleep(5 * time.Second)
	}

	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mainConfig.Host, mainConfig.Port))
	opts.SetClientID(mainConfig.ClientId)
	opts.SetUsername(mainConfig.Username)
	opts.SetPassword(mainConfig.Password)
	mm.Topic = mainConfig.Topic
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	mm.client = mqtt.NewClient(opts)
	if token := mm.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}

}

func (mm *MqttOutEndTarget) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (m *MqttOutEndTarget) OnStreamApproached(data string) error {
	return nil
}
func (mm *MqttOutEndTarget) Stop() {
	mm.client.Disconnect(0)

}
func (mm *MqttOutEndTarget) Reload() {

}
func (mm *MqttOutEndTarget) Pause() {

}
func (mm *MqttOutEndTarget) Status() typex.ResourceState {
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

func (mm *MqttOutEndTarget) Register(outEndId string) error {
	mm.PointId = outEndId
	return nil
}

func (mm *MqttOutEndTarget) Test(outEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
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
func (mm *MqttOutEndTarget) To(data interface{}) error {
	if mm.client != nil {
		return mm.client.Publish(mm.Topic, 2, false, data).Error()
	}
	return errors.New("mqtt client is nil")
}
