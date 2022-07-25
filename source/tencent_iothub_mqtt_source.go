package source

import (
	"encoding/json"
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
// 下行数据
//
type tencentDownMsg struct {
	Method      string      `json:"method"`
	ClientToken string      `json:"clientToken"`
	Params      interface{} `json:"params"`
}

var _PropertyTopic = "$thing/down/property/%v/%v"

// var _EventTopic = "$thing/up/event/%v/%v"
var _ActionTopic = "$thing/down/action/%v/$%v"

//
//
//
type tencentIothubSource struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig common.TencentMqttConfig
}

func NewTencentIothubSource(e typex.RuleX) typex.XSource {
	m := new(tencentIothubSource)
	m.RuleEngine = e
	m.mainConfig = common.TencentMqttConfig{}
	return m
}

/*
*
*
*
 */
func (tc *tencentIothubSource) Init(inEndId string, configMap map[string]interface{}) error {
	tc.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &tc.mainConfig); err != nil {
		return err
	}
	return nil
}

/*
*
*
*
 */
func (tc *tencentIothubSource) Start(cctx typex.CCTX) error {
	tc.Ctx = cctx.Ctx
	tc.CancelCTX = cctx.CancelCTX

	PropertyTopic := fmt.Sprintf(_PropertyTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	// 事件
	// EventTopic := fmt.Sprintf(_PropertyTopic, mainConfig.ProductId, mainConfig.DeviceName)
	// 服务接口
	ActionTopic := fmt.Sprintf(_ActionTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		work, err := tc.RuleEngine.WorkInEnd(tc.RuleEngine.GetInEnd(tc.PointId), string(msg.Payload()))
		if !work {
			glogger.GLogger.Error(err)
		}

	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("Tencent IOTHUB Connected Success")
		client.Subscribe(PropertyTopic, 1, func(c mqtt.Client, m mqtt.Message) {
			msg := tencentDownMsg{}
			if err := json.Unmarshal(m.Payload(), &msg); err != nil {
				glogger.GLogger.Error(err)
			} else {
				glogger.GLogger.Info("Recv: ", m.Topic(), string(m.Payload()))
			}
		})
		client.Subscribe(ActionTopic, 1, func(c mqtt.Client, m mqtt.Message) {
			msg := tencentDownMsg{}
			if err := json.Unmarshal(m.Payload(), &msg); err != nil {
				glogger.GLogger.Error(err)
			} else {
				glogger.GLogger.Info("Recv: ", m.Topic(), string(m.Payload()))
			}
		})
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("Tencent IOTHUB Disconnect: %v, try to reconnect\n", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", tc.mainConfig.Host, tc.mainConfig.Port))
	opts.SetClientID(tc.mainConfig.ClientId)
	opts.SetUsername(tc.mainConfig.Username)
	opts.SetPassword(tc.mainConfig.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	tc.client = mqtt.NewClient(opts)
	if token := tc.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}

}

func (tc *tencentIothubSource) DataModels() []typex.XDataModel {
	return tc.XDataModels
}

func (tc *tencentIothubSource) Stop() {
	tc.client.Disconnect(0)
	tc.CancelCTX()
}
func (tc *tencentIothubSource) Reload() {

}
func (tc *tencentIothubSource) Pause() {

}
func (tc *tencentIothubSource) Status() typex.SourceState {
	if tc.client != nil {
		if tc.client.IsConnected() {
			return typex.SOURCE_UP
		} else {
			return typex.SOURCE_DOWN
		}
	} else {
		return typex.SOURCE_DOWN
	}

}

func (tc *tencentIothubSource) Test(inEndId string) bool {
	if tc.client != nil {
		return tc.client.IsConnected()
	}
	return false
}

func (tc *tencentIothubSource) Enabled() bool {
	return tc.Enable
}
func (tc *tencentIothubSource) Details() *typex.InEnd {
	return tc.RuleEngine.GetInEnd(tc.PointId)
}
func (*tencentIothubSource) Driver() typex.XExternalDriver {
	return nil
}
func (*tencentIothubSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.TENCENT_IOT_HUB, "腾讯云IOTHUB接入支持", common.TencentMqttConfig{})
}

//
// 拓扑
//
func (*tencentIothubSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据
//
func (*tencentIothubSource) DownStream([]byte) {}

//
// 上行数据
//
func (*tencentIothubSource) UpStream() {}
