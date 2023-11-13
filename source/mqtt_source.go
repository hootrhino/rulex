package source

import (
	"fmt"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttInEndSource struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig common.MqttConfig
	status     typex.SourceState
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

	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("Mqtt InEnd Connected Success")
		token := client.Subscribe(mm.mainConfig.SubTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			work, err := mm.RuleEngine.WorkInEnd(mm.RuleEngine.GetInEnd(mm.PointId), string(msg.Payload()))
			if !work {
				glogger.GLogger.Error(err)
			}
		})
		if token.Error() != nil {
			glogger.GLogger.Error(token.Error())
		} else {
			glogger.GLogger.Info("topic:", mm.mainConfig.SubTopic, " subscribed")
		}
		mm.status = typex.SOURCE_UP
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

	opts.SetCleanSession(true)
	opts.SetPingTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	mm.client = mqtt.NewClient(opts)
	if token := mm.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		mm.status = typex.SOURCE_UP
		return nil
	}

}

func (mm *mqttInEndSource) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

func (mm *mqttInEndSource) Stop() {
	mm.status = typex.SOURCE_STOP
	if mm.CancelCTX != nil {
		mm.CancelCTX()
	}
	mm.client.Disconnect(0)
	if mm.client != nil {
		mm.client.Disconnect(0)
		mm.client = nil
	}
}

func (mm *mqttInEndSource) Status() typex.SourceState {
	return mm.status
}

func (mm *mqttInEndSource) Test(inEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
}

func (mm *mqttInEndSource) Details() *typex.InEnd {
	return mm.RuleEngine.GetInEnd(mm.PointId)
}
func (*mqttInEndSource) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*mqttInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*mqttInEndSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*mqttInEndSource) UpStream([]byte) (int, error) {
	return 0, nil
}
