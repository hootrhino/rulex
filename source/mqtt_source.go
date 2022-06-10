package source

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
type mqttInEndSource struct {
	typex.XStatus
	client mqtt.Client
}

func NewMqttInEndSource(inEndId string, e typex.RuleX) typex.XSource {
	m := new(mqttInEndSource)
	m.PointId = inEndId
	m.RuleEngine = e
	return m
}

func (mm *mqttInEndSource) Start(cctx typex.CCTX) error {
	mm.Ctx = cctx.Ctx
	mm.CancelCTX = cctx.CancelCTX

	config := mm.RuleEngine.GetInEnd(mm.PointId).Config
	var mainConfig mqttConfig

	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		work, err := mm.RuleEngine.WorkInEnd(mm.RuleEngine.GetInEnd(mm.PointId), string(msg.Payload()))
		if !work {
			log.Error(err)
		}

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

func (mm *mqttInEndSource) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

func (mm *mqttInEndSource) Stop() {
	mm.client.Disconnect(0)
	mm.CancelCTX()
}
func (mm *mqttInEndSource) Reload() {

}
func (mm *mqttInEndSource) Pause() {

}
func (mm *mqttInEndSource) Status() typex.SourceState {
	if mm.client != nil {
		if mm.client.IsConnected() {
			return typex.SOURCE_UP
		} else {
			return typex.SOURCE_DOWN
		}
	} else {
		return typex.SOURCE_DOWN
	}

}

func (mm *mqttInEndSource) Init(inEndId string, cfg map[string]interface{}) error {
	mm.PointId = inEndId
	return nil
}
func (mm *mqttInEndSource) Test(inEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
}

func (mm *mqttInEndSource) Enabled() bool {
	return mm.Enable
}
func (mm *mqttInEndSource) Details() *typex.InEnd {
	return mm.RuleEngine.GetInEnd(mm.PointId)
}
func (*mqttInEndSource) Driver() typex.XExternalDriver {
	return nil
}
func (*mqttInEndSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.MQTT, "MQTT", mqttConfig{})
}

//
// 拓扑
//
func (*mqttInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
