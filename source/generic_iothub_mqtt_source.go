package source

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

//	{
//	    "method":"${method}_reply",
//	    "requestId":"20a4ccfd-d308",
//	    "code": 0,
//	    "status":"some message"
//	}
const (
	// 属性下发称之为控制指令
	METHOD_CONTROL       string = "control"
	METHOD_CONTROL_REPLY string = "control_reply"
	// 属性下发称之为控制指令[这是为了兼容一些其他平台，实际上和上面功能完全一一样]
	METHOD_PROPERTY       string = "property"
	METHOD_PROPERTY_REPLY string = "property_reply"
	// 动作请求
	METHOD_ACTION       string = "action"
	METHOD_ACTION_REPLY string = "action_reply"
)

const (
	// 属性
	_PropertyTopic      = "$thing/down/property/%v/%v"
	_PropertyUpTopic    = "$thing/up/property/%v/%v"
	_PropertyReplyTopic = "$thing/property/reply/%v/%v"
	// 动作
	_ActionTopic      = "$thing/down/action/%v/%v"
	_ActionUpTopic    = "$thing/up/action/%v/%v"
	_ActionReplyTopic = "$thing/action/reply/%v/%v"
)

/*
*
* 上行数据，包含了上报属性和回复, 用了omitempty属性来灵活处理字段
*
 */
type tencentUpMsg struct {
	Method      string                 `json:"method"`
	ClientToken string                 `json:"clientToken,omitempty"` // 腾讯云
	Id          string                 `json:"id,omitempty"`          // 兼容非腾讯云平台
	Params      map[string]interface{} `json:"params,omitempty"`      // 腾讯云
	Data        map[string]interface{} `json:"data,omitempty"`        // 兼容非腾讯云平台
	Code        int                    `json:"code"`
	Status      string                 `json:"status"`
	// 扩展特性
	UserDefineTopic string `json:"userDefineTopic,omitempty"` // [特性]扩展属性: 用户自定义Topic
	SubDeviceId     string `json:"subDeviceId,omitempty"`     // 网关子设备ID
}

type tencentIothubSource struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig common.TencentMqttConfig
	status     typex.SourceState
}

func NewGenericIothubSource(e typex.RuleX) typex.XSource {
	return NewTencentIothubSource(e)
}
func NewTencentIothubSource(e typex.RuleX) typex.XSource {
	m := new(tencentIothubSource)
	m.RuleEngine = e
	m.mainConfig = common.TencentMqttConfig{}
	m.status = typex.SOURCE_DOWN
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
	ActionTopic := fmt.Sprintf(_ActionTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	// 服务接口
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("Tencent IOTHUB Connected Success")
		tc.subscribe(PropertyTopic)
		tc.subscribe(ActionTopic)
		tc.status = typex.SOURCE_UP
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

func (tc *tencentIothubSource) DataModels() []typex.XDataModel {
	return tc.XDataModels
}

func (tc *tencentIothubSource) Stop() {
	tc.status = typex.SOURCE_STOP
	tc.CancelCTX()
	if tc.client != nil {
		tc.client.Disconnect(0)
		tc.client = nil
	}

}
func (tc *tencentIothubSource) Reload() {

}
func (tc *tencentIothubSource) Pause() {

}
func (tc *tencentIothubSource) Status() typex.SourceState {
	return tc.status
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
	return &typex.XConfig{}
}

// 拓扑
func (*tencentIothubSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据,实际上就是LUA脚本调用的时候写进来的参数
//

func (tc *tencentIothubSource) DownStream(bytes []byte) (int, error) {
	var msg tencentUpMsg
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return 0, err
	}
	//
	var err error
	// 属性回复: 兼容腾讯iothub
	if msg.Method == METHOD_CONTROL_REPLY {
		topic := fmt.Sprintf(_PropertyReplyTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err = tc.client.Publish(topic, 1, false, bytes).Error()
	}
	// 属性回复: 兼容W3C
	if msg.Method == METHOD_PROPERTY_REPLY {
		topic := fmt.Sprintf(_PropertyReplyTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err = tc.client.Publish(topic, 1, false, bytes).Error()
	}
	// 事件调用回复
	if msg.Method == METHOD_ACTION_REPLY {
		topic := fmt.Sprintf(_ActionReplyTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err = tc.client.Publish(topic, 1, false, bytes).Error()
	}
	// 用户自定义Topic
	// iothub:usermsg()
	if msg.Method == "USER_DEFINE_TOPIC" {
		if msg.UserDefineTopic != "" {
			err = tc.client.Publish(msg.UserDefineTopic, 1, false, bytes).Error()
		}
	}
	// 子设备
	// iothub:subdevicemsg()
	if msg.Method == "SUBDEVICE" {
		if msg.UserDefineTopic != "" {
			topic := fmt.Sprintf("$gateway/%s/%s", tc.mainConfig.DeviceName, msg.SubDeviceId)
			err = tc.client.Publish(topic, 1, false, bytes).Error()
		}
	}
	return 0, err
}

// 上行数据
func (*tencentIothubSource) UpStream([]byte) (int, error) {
	return 0, nil
}
func (tc *tencentIothubSource) subscribe(topic string) {
	token := tc.client.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
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
