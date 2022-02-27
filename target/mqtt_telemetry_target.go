package target

/*
*
* 这个是一个比较特殊的资源出口, 专门用来遥测上报数据, 实现和任意 Mqtt Broker的交互接口
*
 */
import (
	"encoding/json"
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
type mqttTelemetryConfig struct {
	Host         string `json:"host" validate:"required"`
	Port         int    `json:"port" validate:"required"`
	S2CTopic     string `json:"S2CTopic" validate:"required"`     // 这个Topic是专门留给服务器下发指令用的
	ToplogyTopic string `json:"toplogyTopic" validate:"required"` // 定时上报拓扑结构的 Topic
	LogTopic     string `json:"logTopic" validate:"required"`     // 定时上报拓扑结构的 Topic
	DataTopic    string `json:"dataTopic" validate:"required"`    // 自定义上报数据的 Topic
	StateTopic   string `json:"stateTopic" validate:"required"`   // 定时上报状态的 Topic
	ClientId     string `json:"clientId" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

//
type MqttTelemetryTarget struct {
	typex.XStatus
	client       mqtt.Client
	S2CTopic     string `json:"S2CTopic" validate:"required"`     // 这个Topic是专门留给服务器下发指令用的
	ToplogyTopic string `json:"toplogyTopic" validate:"required"` // 定时上报拓扑结构的 Topic
	DataTopic    string `json:"dataTopic" validate:"required"`    // 上报数据的 Topic
	StateTopic   string `json:"stateTopic" validate:"required"`   // 定时上报状态的 Topic
}

//
// 服务器消息
//
type s2cCommand struct {
	Cmd  string        `json:"cmd" validate:"required"`
	Args []interface{} `json:"args" validate:"required"`
}
type c2sCommand struct {
	Type   string
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
func (mm *MqttTelemetryTarget) Start(cctx typex.CCTX) error {
	mm.Ctx = cctx.Ctx
	mm.CancelCTX = cctx.CancelCTX
	outEnd := mm.RuleEngine.GetOutEnd(mm.PointId)
	config := outEnd.Config
	var mainConfig mqttTelemetryConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
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
				// 内置的几个简单指令, 后期会扩展
				if cmd.Cmd == "get-state" {
					token := mm.client.Publish(mainConfig.StateTopic, 0, false, c2sCommand{
						Type:   cmd.Cmd,
						Result: "running",
					})
					if token.Error() != nil {
						log.Error(token.Error())
					}
				} else if cmd.Cmd == "get-topology" {
					token := mm.client.Publish(mainConfig.ToplogyTopic, 0, false, c2sCommand{
						Type:   cmd.Cmd,
						Result: []string{},
					})
					if token.Error() != nil {
						log.Error(token.Error())
					}

				} else if cmd.Cmd == "get-log" {
					//
					//  {"cmd" : "get-log", args: [1, 100]}
					//
					var offset float64
					var size float64
					if len(cmd.Args) == 2 {
						offset = cmd.Args[0].(float64)
						size = cmd.Args[1].(float64)
					} else {
						offset = 0
						size = 20
					}
					c := c2sCommand{
						Type:   cmd.Cmd,
						Result: core.GLOBAL_LOGGER.Slot()[int(offset):int(size)],
					}
					bytes, _ := json.Marshal(c)
					token := mm.client.Publish(mainConfig.LogTopic, 0, false, string(bytes))
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
	//
	mm.DataTopic = mainConfig.DataTopic
	mm.S2CTopic = mainConfig.S2CTopic
	mm.StateTopic = mainConfig.StateTopic
	mm.ToplogyTopic = mainConfig.ToplogyTopic
	//
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

/*
*
* 物模型接口定义
*
 */
func (mm *MqttTelemetryTarget) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (mm *MqttTelemetryTarget) OnStreamApproached(data string) error {
	//-----------------------------------------------------------------------------------
	// 对于遥测组件来说, 所有进来的 [非规则引擎数据] 数据全发到 State topic, 传到远程处理.
	//-----------------------------------------------------------------------------------
	if mm.client != nil {
		return mm.client.Publish(mm.StateTopic, 2, false, data).Error()
	}
	return nil
}
func (mm *MqttTelemetryTarget) Stop() {
	mm.client.Disconnect(0)
	mm.CancelCTX()

}
func (mm *MqttTelemetryTarget) Reload() {

}
func (mm *MqttTelemetryTarget) Pause() {

}
func (mm *MqttTelemetryTarget) Status() typex.SourceState {
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
func (mm *MqttTelemetryTarget) Init(outEndId string, cfg map[string]interface{}) error {
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

/*
*
* 配置
*
 */
func (*MqttTelemetryTarget) Configs() *typex.XConfig {
	return core.GenOutConfig(typex.MQTT_TELEMETRY_TARGET, "MQTT_TELEMETRY_TARGET", httpConfig{})
}
