package source

import (
	"fmt"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _PropertyTopic = "$thing/down/property/%v/%v"
var _EventTopic = "$thing/up/event/%v/%v"
var _ActionTopic = "$thing/down/action/%v/$%v"

//
//
//
type tencentMqttConfig struct {
	ProductId  string `json:"productId" validate:"required" title:"产品名" info:""`
	DeviceName string `json:"deviceName" validate:"required" title:"设备名" info:""`
	//
	Host string `json:"host" validate:"required" title:"服务地址" info:""`
	Port int    `json:"port" validate:"required" title:"服务端口" info:""`
	//
	ClientId string `json:"clientId" validate:"required" title:"客户端ID" info:""`
	Username string `json:"username" validate:"required" title:"连接账户" info:""`
	Password string `json:"password" validate:"required" title:"连接密码" info:""`
}

//
//
//
type tencentIothubSource struct {
	typex.XStatus
	client mqtt.Client
}

func NewtencentIothubSource(e typex.RuleX) typex.XSource {
	m := new(tencentIothubSource)
	m.RuleEngine = e
	return m
}

//
//
//
func (tc *tencentIothubSource) Start(cctx typex.CCTX) error {
	tc.Ctx = cctx.Ctx
	tc.CancelCTX = cctx.CancelCTX

	config := tc.RuleEngine.GetInEnd(tc.PointId).Config
	var mainConfig tencentMqttConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	PropertyTopic := fmt.Sprintf(_PropertyTopic, mainConfig.ProductId, mainConfig.DeviceName)
	// 事件
	// EventTopic := fmt.Sprintf(_PropertyTopic, mainConfig.ProductId, mainConfig.DeviceName)
	// 服务接口
	ActionTopic := fmt.Sprintf(_ActionTopic, mainConfig.ProductId, mainConfig.DeviceName)

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
			glogger.GLogger.Debug("Recv: ", m.Topic(), m.Payload())
		})
		client.Subscribe(ActionTopic, 1, func(c mqtt.Client, m mqtt.Message) {
			glogger.GLogger.Debug("Recv: ", m.Topic(), m.Payload())
		})
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("Tencent IOTHUB Disconnect: %v, try to reconnect\n", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mainConfig.Host, mainConfig.Port))
	opts.SetClientID(mainConfig.ClientId)
	opts.SetUsername(mainConfig.Username)
	opts.SetPassword(mainConfig.Password)
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

func (tc *tencentIothubSource) Init(inEndId string, cfg map[string]interface{}) error {
	tc.PointId = inEndId
	return nil
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
	return core.GenInConfig(typex.TENCENT_IOT_HUB, "腾讯云IOTHUB接入支持", mqttConfig{})
}

//
// 拓扑
//
func (*tencentIothubSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
