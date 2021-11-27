package target

/*
*
* 用来遥测上报数据
*
 */
import (
	"encoding/json"
	"errors"
	"fmt"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
type mqttTelemetryConfig struct {
	Host         string `json:"host" validate:"required"`
	Port         int    `json:"port" validate:"required"`
	S2CTopic     string `json:"S2CTopic" validate:"required"`     // 这个Topic是专门留给服务器下发指令用的
	ToplogyTopic string `json:"toplogyTopic" validate:"required"` // 定时上报拓扑结构的 Topic
	DataTopic    string `json:"dataTopic" validate:"required"`    // 上报数据的 Topic
	StateTopic   string `json:"stateTopic" validate:"required"`   // 定时上报状态的 Topic
	ClientId     string `json:"clientId" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

//
type MqttTelemetryTarget struct {
	typex.XStatus
	client    mqtt.Client
	DataTopic string
}

//
// 服务器消息
//
type s2cCommand struct {
	Cmd  string
	Args []string
}
type c2sCommand struct {
	Result interface{}
}

func NewMqttTelemetryTarget(e typex.RuleX) typex.XTarget {
	m := new(MqttTelemetryTarget)
	m.RuleEngine = e
	return m
}
func (*MqttTelemetryTarget) Driver() typex.XExternalDriver {
	return nil
}
func (mm *MqttTelemetryTarget) Start() error {
	outEnd := mm.RuleEngine.GetOutEnd(mm.PointId)
	config := outEnd.Config
	var mainConfig mqttTelemetryConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Mqtt OutEnd Connected Success")
		client.Subscribe(mainConfig.S2CTopic, 2, func(c mqtt.Client, m mqtt.Message) {
			var cmd s2cCommand
			if err := json.Unmarshal(m.Payload(), &cmd); err != nil {
				log.Error(err)
			} else {
				if cmd.Cmd == "get-state" {
					token := mm.client.Publish(mainConfig.StateTopic, 0, false, c2sCommand{
						Result: "running",
					})
					if token.Error() != nil {
						log.Error(token.Error())
					}
				} else if cmd.Cmd == "get-toplogy" {
					token := mm.client.Publish(mainConfig.ToplogyTopic, 0, false, c2sCommand{
						Result: []string{},
					})
					if token.Error() != nil {
						log.Error(token.Error())
					}

				} else if cmd.Cmd == "get-log" {
					token := mm.client.Publish(mainConfig.ToplogyTopic, 0, false, c2sCommand{
						Result: []string{},
					})
					if token.Error() != nil {
						log.Error(token.Error())
					}
				} else {
					log.Error("Unsupported command:" + cmd.Cmd)
				}
			}
		})
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

func (mm *MqttTelemetryTarget) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (m *MqttTelemetryTarget) OnStreamApproached(data string) error {
	return nil
}
func (mm *MqttTelemetryTarget) Stop() {
	mm.client.Disconnect(0)

}
func (mm *MqttTelemetryTarget) Reload() {

}
func (mm *MqttTelemetryTarget) Pause() {

}
func (mm *MqttTelemetryTarget) Status() typex.ResourceState {
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

func (mm *MqttTelemetryTarget) Register(outEndId string) error {
	mm.PointId = outEndId
	return nil
}

func (mm *MqttTelemetryTarget) Test(outEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
}

func (mm *MqttTelemetryTarget) Enabled() bool {
	return mm.Enable
}
func (mm *MqttTelemetryTarget) Details() *typex.OutEnd {
	return mm.RuleEngine.GetOutEnd(mm.PointId)
}

//
//
//
func (mm *MqttTelemetryTarget) To(data interface{}) error {
	if mm.client != nil {
		return mm.client.Publish(mm.DataTopic, 2, false, data).Error()
	}
	return errors.New("mqtt client is nil")
}
