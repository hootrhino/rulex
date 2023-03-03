package source

/*
*
* ithings 支持, 其本质上是个MQTT客户端, 和ithings进行交互
*
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	// 属性
	_ithings_PropertyTopic   = "$thing/down/property/%v/%v"
	_ithings_PropertyUpTopic = "$thing/up/property/%v/%v"
	// 动作
	_ithings_ActionTopic   = "$thing/down/action/%v/%v"
	_ithings_ActionUpTopic = "$thing/up/action/%v/%v"
	// 事件
	_ithings_EventUpTopic = "$thing/up/event/%v/%v"
)

/*
*
* 上行数据，包含了上报属性和回复, 用了omitempty属性来灵活处理字段
*
 */
type ithingsUpMsg struct {
	Method      string                 `json:"method"`
	ClientToken string                 `json:"clientToken,omitempty"`
	ActionId    string                 `json:"actionId,omitempty"`
	EventId     string                 `json:"eventId,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Code        int                    `json:"code"`
	Status      string                 `json:"status"`
	// 扩展特性
	UserDefineTopic string `json:"userDefineTopic,omitempty"` // [特性]扩展属性: 用户自定义Topic
	SubDeviceId     string `json:"subDeviceId,omitempty"`     // 网关子设备ID
}

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
	tc.status = typex.SOURCE_STOP
	tc.CancelCTX()
	if tc.client != nil {
		tc.client.Disconnect(0)
		tc.client = nil
	}
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
	return &typex.XConfig{}
}

// 拓扑
func (*ithings) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据,实际上就是LUA脚本调用的时候写进来的参数
//

func (tc *ithings) DownStream(bytes []byte) (int, error) {
	msg := ithingsUpMsg{}
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return 0, err
	}
	// 属性上报数据有2类：
	// 1 网关本身
	// 2 子设备
	methods := map[string]string{
		"report":         _ithings_PropertyUpTopic, // 属性上报
		"controlReply":   _ithings_PropertyUpTopic, // 属性下发结果上报
		"getStatusReply": _ithings_PropertyUpTopic, // 获取实时状态上报
		"eventPost":      _ithings_EventUpTopic,    // 事件上报
	}

	if methods[msg.Method] != "" {
		topic := fmt.Sprintf(methods[msg.Method], tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err := tc.client.Publish(topic, 1, false, bytes).Error()
		if err != nil {
			glogger.GLogger.Error(err)
		}
	} else {
		glogger.GLogger.Error(errors.New("unsupported method:" + msg.Method))
	}

	return 0, nil
}

// 上行数据
func (*ithings) UpStream([]byte) (int, error) {
	return 0, nil
}

func (tc *ithings) subscribe(topic string) {
	token := tc.client.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
		// 所有的消息都给丢进规则引擎里面, 交给用户的lua脚本来处理
		work, err := tc.RuleEngine.WorkInEnd(tc.RuleEngine.GetInEnd(tc.PointId), string(msg.Payload()))
		if !work {
			glogger.GLogger.Error(err)
		}
	})
	if token.Error() != nil {
		glogger.GLogger.Error(token.Error())
	} else {
		glogger.GLogger.Info("topic:", topic, " subscribed")
	}

}
