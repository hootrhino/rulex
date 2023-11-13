package source

/*
*
* iothub 支持, 其本质上是个MQTT客户端, 和iothub进行交互
*
 */

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	// 属性
	_iothub_PropertyTopic   = "$thing/down/property/%v/%v"
	_iothub_PropertyUpTopic = "$thing/up/property/%v/%v"
	// 动作
	_iothub_ActionTopic   = "$thing/down/action/%v/%v"
	_iothub_ActionUpTopic = "$thing/up/action/%v/%v"
	// 子设备拓扑 [当前版本iothub不支持]
	_iothub_subTopologyDown = "$thing/down/operation/%v/%v"
	_iothub_subTopologyUp   = "$thing/up/operation/%v/%v"
)

// 来自外面的数据,实际上就是LUA脚本调用的时候写进来的参数
const (
	// 属性下发称之为控制指令
	_iothub_METHOD_CONTROL       string = "control"
	_iothub_METHOD_CONTROL_REPLY string = "control_reply"
	// 属性下发称之为控制指令[这是为了兼容一些其他平台，实际上和上面功能完全一一样]
	_iothub_METHOD_PROPERTY       string = "property"
	_iothub_METHOD_PROPERTY_REPLY string = "property_reply"
	// 动作请求
	_iothub_METHOD_ACTION       string = "action"
	_iothub_METHOD_ACTION_REPLY string = "action_reply"
)

/*
*
* 上行数据，包含了上报属性和回复, 用了omitempty属性来灵活处理字段
*
 */
type iothubUpMsg struct {
	Method      string `json:"method"`
	ClientToken string `json:"clientToken,omitempty"`
	ActionId    string `json:"actionId,omitempty"`
	EventId     string `json:"eventId,omitempty"`
	// 时间戳
	Timestamp int64                  `json:"timestamp,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"` // 上报
	Data      map[string]interface{} `json:"data,omitempty"`   // 下发
	Code      int                    `json:"code"`
	Status    string                 `json:"status"`
	//------------------------------------------------------------------
	// 扩展特性
	//------------------------------------------------------------------
	UserDefineTopic string `json:"userDefineTopic,omitempty"` // [特性]扩展属性: 用户自定义Topic
	SubDeviceId     string `json:"subDeviceId,omitempty"`     // 网关子设备ID

}

func (it iothubUpMsg) String() string {
	b, _ := json.Marshal(it)
	return string(b)
}

/*
*
* 下发拓扑结构
*
 */
// {
// 	"method":"describesubDevices",
// 	"clientToken":"3160be0b-6d4f-e6fa-d614-8fb422c0d16c",
// 	"timestamp":1681625336600,
// 	"status":"成功",
// 	"payload":{
// 		"devices":[
// 			{
// 				"productID":"268dGhSTdOE",
// 				"deviceName":"RULEX-大屏1"
// 			}
// 		]
// 	}
// }
type TopologyDevice struct {
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

// $thing/status/${productid}/${devicename}
func (tt TopologyDevice) OnlineTopic() string {
	return fmt.Sprintf("$thing/status/%s/%s", tt.ProductID, tt.DeviceName)
}

// $thing/status/${productid}/${devicename}
func (tt TopologyDevice) OfflineTopic() string {
	return fmt.Sprintf("$thing/status/%s/%s", tt.ProductID, tt.DeviceName)
}

// $thing/up/property/{ProductID}/{DeviceName}
func (tt TopologyDevice) ReportTopic() string {
	return fmt.Sprintf("$thing/up/property/%s/%s", tt.ProductID, tt.DeviceName)
}

type TopologyPayload struct {
	Devices []TopologyDevice `json:"devices"`
}
type TopologyDownMsg struct {
	Method      string          `json:"method"`
	ClientToken string          `json:"clientToken"`
	Timestamp   int64           `json:"timestamp"`
	Status      string          `json:"status"`
	Payload     TopologyPayload `json:"payload"`
}

type iothub struct {
	typex.XStatus
	client        mqtt.Client
	mainConfig    common.GenericIoTHUBMqttConfig
	status        typex.SourceState
	subDevices    map[string]TopologyDevice // 子设备
	topologyReady bool
	// Topic
	PropertyUpTopic   string
	PropertyDownTopic string
	ActionUpTopic     string
	ActionDownTopic   string
	TopologyTopicDown string
}

func NewIoTHubSource(e typex.RuleX) typex.XSource {
	m := new(iothub)
	m.RuleEngine = e
	m.mainConfig = common.GenericIoTHUBMqttConfig{
		Mode: "DC",
	}
	m.subDevices = map[string]TopologyDevice{}
	m.topologyReady = false
	m.status = typex.SOURCE_DOWN
	return m
}

/*
*
*
*
 */
func (tc *iothub) Init(inEndId string, configMap map[string]interface{}) error {
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
func (tc *iothub) Start(cctx typex.CCTX) error {
	tc.Ctx = cctx.Ctx
	tc.CancelCTX = cctx.CancelCTX
	// 服务接口
	//
	tc.PropertyDownTopic = fmt.Sprintf(_iothub_PropertyTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	tc.PropertyUpTopic = fmt.Sprintf(_iothub_PropertyUpTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	tc.ActionDownTopic = fmt.Sprintf(_iothub_ActionTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	tc.ActionUpTopic = fmt.Sprintf(_iothub_ActionUpTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	// 网关模式
	if tc.mainConfig.Mode == "GW" {
		glogger.GLogger.Info("Connect iothub with Gateway Mode")
		TopologyTopicDown := fmt.Sprintf(_iothub_subTopologyDown, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		TopologyTopicUp := fmt.Sprintf(_iothub_subTopologyUp, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		token3 := tc.client.Subscribe(TopologyTopicDown, 1, func(c mqtt.Client, m mqtt.Message) {
			topologyDownMsg := TopologyDownMsg{
				Payload: TopologyPayload{
					Devices: []TopologyDevice{},
				},
			}
			if err := json.Unmarshal(m.Payload(), &topologyDownMsg); err != nil {
				glogger.GLogger.Error(err)
				return
			}
			glogger.GLogger.Info("Topology fetch success")
			for _, dev := range topologyDownMsg.Payload.Devices {
				tc.subDevices[dev.DeviceName] = dev
			}
			tc.topologyReady = true
		})
		if token3.Error() != nil {
			glogger.GLogger.Info("Topology fetch error:", token3.Error())
			return token3.Error()
		}
		// 请求下发拓扑
		token4 := tc.client.Publish(TopologyTopicUp, 1, false, `{"method":"describeSubDevices"}`)
		if token4.Error() != nil {
			glogger.GLogger.Error("Topology Publish error:", token3.Error())
			return token4.Error()
		}
	}
	if tc.mainConfig.Mode == "DC" {
		glogger.GLogger.Info("Connect iothub with Direct Connect Mode")
	}
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("IOTHUB Connected Success")
		if err := tc.subscribe(tc.PropertyDownTopic); err != nil {
			glogger.GLogger.Error(err)
		}
		if err := tc.subscribe(tc.ActionDownTopic); err != nil {
			glogger.GLogger.Error(err)
		}
		if tc.mainConfig.Mode == "GW" {
			if err := tc.subscribe(tc.TopologyTopicDown); err != nil {
				glogger.GLogger.Error(err)
			}
		}

	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("IOTHUB Disconnect: %v, %v try to reconnect", err, tc.status)
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
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetAutoReconnect(false)
	opts.SetMaxReconnectInterval(0)
	tc.client = mqtt.NewClient(opts)
	if token := tc.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	tc.status = typex.SOURCE_UP
	return nil

}

func (tc *iothub) DataModels() []typex.XDataModel {
	return tc.XDataModels
}

func (tc *iothub) Stop() {
	tc.status = typex.SOURCE_DOWN
	if tc.CancelCTX != nil {
		tc.CancelCTX()
	}
	if tc.client != nil {
		tc.client.Unsubscribe(tc.PropertyDownTopic)
		tc.client.Unsubscribe(tc.ActionDownTopic)
		tc.client.Unsubscribe(tc.TopologyTopicDown)
		tc.client.Disconnect(100)
	}
}

func (tc *iothub) Status() typex.SourceState {
	if tc.client != nil {
		if tc.client.IsConnectionOpen() {
			return typex.SOURCE_UP
		}
	}
	return typex.SOURCE_DOWN
}

func (tc *iothub) Test(inEndId string) bool {
	if tc.client != nil {
		return tc.client.IsConnected()
	}
	return false
}


func (tc *iothub) Details() *typex.InEnd {
	return tc.RuleEngine.GetInEnd(tc.PointId)
}
func (*iothub) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*iothub) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

func (tc *iothub) DownStream(bytes []byte) (int, error) {
	var msg iothubUpMsg
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return 0, err
	}
	//
	var err error
	// 属性回复: 兼容腾讯iothub
	if msg.Method == _iothub_METHOD_CONTROL_REPLY {
		topic := fmt.Sprintf(tc.PropertyUpTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err = tc.client.Publish(topic, 1, false, bytes).Error()
	}
	// 属性回复: 兼容W3C
	if msg.Method == _iothub_METHOD_PROPERTY_REPLY {
		topic := fmt.Sprintf(tc.ActionUpTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err = tc.client.Publish(topic, 1, false, bytes).Error()
	}
	// 事件调用回复
	if msg.Method == _iothub_METHOD_ACTION_REPLY {
		topic := fmt.Sprintf(tc.ActionUpTopic, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
		err = tc.client.Publish(topic, 1, false, bytes).Error()
	}
	// 用户自定义Topic
	// iothub:usermsg()
	if msg.Method == "user_topic" {
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
func (*iothub) UpStream([]byte) (int, error) {
	return 0, nil
}

func (tc *iothub) subscribe(topic string) error {
	token := tc.client.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
		// 所有的消息都给丢进规则引擎里面, 交给用户的lua脚本来处理
		work, err := tc.RuleEngine.WorkInEnd(tc.RuleEngine.GetInEnd(tc.PointId), string(msg.Payload()))
		if !work {
			glogger.GLogger.Error(err)
		}
	})
	if token.Error() != nil {
		glogger.GLogger.Error(token.Error())
		return token.Error()
	} else {
		glogger.GLogger.Info("topic:", topic, " subscribed")
		return nil
	}

}

// --------------------------------------------------------------------------------------------------
// User data
// --------------------------------------------------------------------------------------------------

/*
*
* 下发拓扑结构
*
 */
// {
// 	"method":"describesubDevices",
// 	"clientToken":"3160be0b-6d4f-e6fa-d614-8fb422c0d16c",
// 	"timestamp":1681625336600,
// 	"status":"成功",
// 	"payload":{
// 		"devices":[
// 			{
// 				"productID":"268dGhSTdOE",
// 				"deviceName":"RULEX-大屏1"
// 			}
// 		]
// 	}
// }
type IotHUBTopologyDevice struct {
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

// $thing/status/${productid}/${devicename}
func (tt IotHUBTopologyDevice) OnlineTopic() string {
	return fmt.Sprintf("$thing/status/%s/%s", tt.ProductID, tt.DeviceName)
}

// $thing/status/${productid}/${devicename}
func (tt IotHUBTopologyDevice) OfflineTopic() string {
	return fmt.Sprintf("$thing/status/%s/%s", tt.ProductID, tt.DeviceName)
}

// $thing/up/property/{ProductID}/{DeviceName}
func (tt IotHUBTopologyDevice) ReportTopic() string {
	return fmt.Sprintf("$thing/up/property/%s/%s", tt.ProductID, tt.DeviceName)
}

type IotHUBTopologyPayload struct {
	Devices []IotHUBTopologyDevice `json:"devices"`
}
type IotHUBTopologyDownMsg struct {
	Method      string          `json:"method"`
	ClientToken string          `json:"clientToken"`
	Timestamp   int64           `json:"timestamp"`
	Status      string          `json:"status"`
	Payload     TopologyPayload `json:"payload"`
}
