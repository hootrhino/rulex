package target

import (
	"errors"
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
	Host      string `json:"host" validate:"required"`
	Port      int    `json:"port" validate:"required"`
	DataTopic string `json:"dataTopic" validate:"required"` // 上报数据的 Topic
	ClientId  string `json:"clientId" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

//
type mqttOutEndTarget struct {
	typex.XStatus
	client    mqtt.Client
	DataTopic string
}

func NewMqttTarget(e typex.RuleX) typex.XTarget {
	m := new(mqttOutEndTarget)
	m.RuleEngine = e
	return m
}
func (*mqttOutEndTarget) Driver() typex.XExternalDriver {
	return nil
}
func (mm *mqttOutEndTarget) Start() error {
	outEnd := mm.RuleEngine.GetOutEnd(mm.PointId)
	config := outEnd.Config
	var mainConfig mqttConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Mqtt OutEnd Connected Success")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Warnf("Connect lost: %v, try to reconnect\n", err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mainConfig.Host, mainConfig.Port))
	opts.SetClientID(mainConfig.ClientId)
	opts.SetUsername(mainConfig.Username)
	opts.SetPassword(mainConfig.Password)
	mm.DataTopic = mainConfig.DataTopic
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetPingTimeout(3 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	mm.client = mqtt.NewClient(opts)
	token := mm.client.Connect()
	token.WaitTimeout(3 * time.Second)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}

}

func (mm *mqttOutEndTarget) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (m *mqttOutEndTarget) OnStreamApproached(data string) error {
	return nil
}
func (mm *mqttOutEndTarget) Stop() {
	mm.client.Disconnect(0)

}
func (mm *mqttOutEndTarget) Reload() {

}
func (mm *mqttOutEndTarget) Pause() {

}
func (mm *mqttOutEndTarget) Status() typex.SourceState {
	if mm.client != nil {
		if mm.client.IsConnectionOpen() {
			return typex.UP
		} else {
			return typex.DOWN
		}
	} else {
		return typex.DOWN
	}

}

func (mm *mqttOutEndTarget) Register(outEndId string) error {
	mm.PointId = outEndId
	return nil
}

func (mm *mqttOutEndTarget) Test(outEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
}

func (mm *mqttOutEndTarget) Enabled() bool {
	return mm.Enable
}
func (mm *mqttOutEndTarget) Details() *typex.OutEnd {
	return mm.RuleEngine.GetOutEnd(mm.PointId)
}

//
//
//
func (mm *mqttOutEndTarget) To(data interface{}) error {
	if mm.client != nil {
		return mm.client.Publish(mm.DataTopic, 2, false, data).Error()
	}
	return errors.New("mqtt client is nil")
}

/*
*
* 配置
*
 */
func (*mqttOutEndTarget) Configs() *typex.XConfig {
	return core.GenOutConfig(typex.MQTT_TARGET, "MQTT_TARGET", httpConfig{})
}
