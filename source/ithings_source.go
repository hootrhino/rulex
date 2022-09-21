package source

/*
*
* ithings 支持, 其本质上是个MQTT客户端, 和ithings进行交互
*
 */

import (
	"fmt"
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	// 属性下发称之为控制指令
	ithings_METHOD_CONTROL       string = "control"
	ithings_METHOD_CONTROL_REPLY string = "control_reply"
	// 属性下发称之为控制指令[这是为了兼容一些其他平台，实际上和上面功能完全一一样]
	ithings_METHOD_PROPERTY       string = "property"
	ithings_METHOD_PROPERTY_REPLY string = "property_reply"
	// 动作请求
	ithings_METHOD_ACTION       string = "action"
	ithings_METHOD_ACTION_REPLY string = "action_reply"
)

const (
	// 属性
	_ithings_PropertyTopic      = "$thing/down/property/%v/%v"
	_ithings_PropertyUpTopic    = "$thing/up/property/%v/%v"
	_ithings_PropertyReplyTopic = "$thing/property/reply/%v/%v"
	// 动作
	_ithings_ActionTopic      = "$thing/down/action/%v/%v"
	_ithings_ActionUpTopic    = "$thing/up/action/%v/%v"
	_ithings_ActionReplyTopic = "$thing/action/reply/%v/%v"
)

/*
*
* 上行数据，包含了上报属性和回复, 用了omitempty属性来灵活处理字段
*
 */
type ithingsUpMsg struct {
	Method      string                 `json:"method"`
	ClientToken string                 `json:"clientToken,omitempty"` // 腾讯云
	Id          string                 `json:"id,omitempty"`          // 兼容非腾讯云平台
	Params      map[string]interface{} `json:"params,omitempty"`      // 腾讯云
	Data        map[string]interface{} `json:"data,omitempty"`        // 兼容非腾讯云平台
	Code        int                    `json:"code"`
	Status      string                 `json:"status"`
}

//
//
//
type ithings struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig common.IThingsMqttConfig
	status     typex.SourceState
}

func NewIThingsSource(e typex.RuleX) typex.XSource {
	m := new(ithings)
	m.RuleEngine = e
	m.mainConfig = common.IThingsMqttConfig{}
	m.status = typex.SOURCE_DOWN
	return m
}

/*
*
*
*
 */
func (tc *ithings) Init(inEndId string, configMap map[string]interface{}) error {
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
func (tc *ithings) Start(cctx typex.CCTX) error {
	tc.Ctx = cctx.Ctx
	tc.CancelCTX = cctx.CancelCTX

	PropertyTopic := fmt.Sprintf(_ithings_PropertyTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	// 事件
	ActionTopic := fmt.Sprintf(_ithings_ActionTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	// 服务接口
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("IThings IOTHUB Connected Success")
		tc.subscribe(PropertyTopic)
		tc.subscribe(ActionTopic)
		tc.status = typex.SOURCE_UP
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("IThings IOTHUB Disconnect: %v, try to reconnect\n", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", tc.mainConfig.Host, tc.mainConfig.Port))
	opts.SetClientID(tc.mainConfig.ClientId)
	opts.SetUsername(tc.mainConfig.Username)
	opts.SetPassword(tc.mainConfig.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetCleanSession(true)
	opts.SetPingTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	tc.client = mqtt.NewClient(opts)
	if token := tc.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil

}

func (tc *ithings) DataModels() []typex.XDataModel {
	return tc.XDataModels
}

func (tc *ithings) Stop() {
	tc.client.Disconnect(0)
	tc.CancelCTX()
	tc.status = typex.SOURCE_STOP
}
func (tc *ithings) Reload() {

}
func (tc *ithings) Pause() {

}
func (tc *ithings) Status() typex.SourceState {
	return tc.status
}

func (tc *ithings) Test(inEndId string) bool {
	if tc.client != nil {
		return tc.client.IsConnected()
	}
	return false
}

func (tc *ithings) Enabled() bool {
	return tc.Enable
}
func (tc *ithings) Details() *typex.InEnd {
	return tc.RuleEngine.GetInEnd(tc.PointId)
}
func (*ithings) Driver() typex.XExternalDriver {
	return nil
}
func (*ithings) Configs() *typex.XConfig {
	return core.GenInConfig(typex.ITHINGS_IOT_HUB, "ITHINGS IOTHUB接入支持", common.IThingsMqttConfig{})
}

//
// 拓扑
//
func (*ithings) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据,实际上就是LUA脚本调用的时候写进来的参数
//

func (tc *ithings) DownStream(bytes []byte) (int, error) {
	return 0, nil
}

//
// 上行数据
//
func (*ithings) UpStream([]byte) (int, error) {
	return 0, nil
}

func (tc *ithings) subscribe(topic string) {
	token := tc.client.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
		// Ithings的消息不走规则引擎, 它纯粹是一个外部资源
	})
	if token.Error() != nil {
		glogger.GLogger.Error(token.Error())
	} else {
		glogger.GLogger.Info("topic:", topic, " subscribed")
	}

}
