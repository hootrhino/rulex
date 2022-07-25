package source

import (
	"fmt"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
type mqttInEndSource struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig common.MqttConfig
}

func NewMqttInEndSource(e typex.RuleX) typex.XSource {
	m := new(mqttInEndSource)
	m.RuleEngine = e
	return m
}

func (mm *mqttInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	mm.PointId = inEndId

	if err := utils.BindSourceConfig(configMap, &mm.mainConfig); err != nil {
		return err
	}
	return nil
}

/*
*
* Start
*
 */
func (mm *mqttInEndSource) Start(cctx typex.CCTX) error {
	mm.Ctx = cctx.Ctx
	mm.CancelCTX = cctx.CancelCTX

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		work, err := mm.RuleEngine.WorkInEnd(mm.RuleEngine.GetInEnd(mm.PointId), string(msg.Payload()))
		if !work {
			glogger.GLogger.Error(err)
		}

	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("Mqtt InEnd Connected Success")
		client.Subscribe(mm.mainConfig.SubTopic, 1, nil)
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("Connect lost: %v, try to reconnect\n", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mm.mainConfig.Host, mm.mainConfig.Port))
	opts.SetClientID(mm.mainConfig.ClientId)
	opts.SetUsername(mm.mainConfig.Username)
	opts.SetPassword(mm.mainConfig.Password)
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
	return core.GenInConfig(typex.MQTT, "MQTT", common.MqttConfig{})
}

//
// 拓扑
//
func (*mqttInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据
//
func (*mqttInEndSource) DownStream([]byte) {}

//
// 上行数据
//
func (*mqttInEndSource) UpStream() {}
