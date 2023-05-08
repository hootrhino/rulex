package source

/*
*
* ithings 支持, 其本质上是个MQTT客户端, 和ithings进行交互
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
	_ithings_PropertyTopic   = "$thing/down/property/%v/%v"
	_ithings_PropertyUpTopic = "$thing/up/property/%v/%v"
	// 动作
	_ithings_ActionTopic   = "$thing/down/action/%v/%v"
	_ithings_ActionUpTopic = "$thing/up/action/%v/%v"
	// 子设备拓扑
	_ithings_subTopologyDown = "$gateway/down/operation/%v/%v"
	_ithings_subTopologyUp   = "$gateway/up/operation/%v/%v"
)

/*
*
* 上行数据，包含了上报属性和回复, 用了omitempty属性来灵活处理字段
*
 */
type ithingsUpMsg struct {
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
	SubDeviceName string `json:"subDeviceName,omitempty"` // 网关子设备ID
}

func (it ithingsUpMsg) string() string {
	b, _ := json.Marshal(it)
	return string(b)
}

type ithings struct {
	typex.XStatus
	client        mqtt.Client
	mainConfig    common.IThingsMqttConfig
	status        typex.SourceState
	subDevices    map[string]TopologyDevice // 子设备
	topologyReady bool
}

func NewIThingsSource(e typex.RuleX) typex.XSource {
	m := new(ithings)
	m.RuleEngine = e
	m.mainConfig = common.IThingsMqttConfig{}
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
	TopologyTopicDown := fmt.Sprintf(_ithings_subTopologyDown, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	TopologyTopicUp := fmt.Sprintf(_ithings_subTopologyUp, tc.mainConfig.ProductId, tc.mainConfig.DeviceName)
	// 服务接口
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("IThings IOTHUB Connected Success")
		tc.subscribe(PropertyTopic)
		tc.subscribe(ActionTopic)

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

		} else {
			// 请求下发拓扑
			token4 := tc.client.Publish(TopologyTopicUp, 1, false, `{"method":"describeSubDevices"}`)
			if token4.Error() != nil {
				glogger.GLogger.Error("Topology Publish error:", token3.Error())

			}
		}

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
	// methods
	// 	"report":         _ithings_PropertyUpTopic, // 属性上报
	// 	"controlReply":   _ithings_PropertyUpTopic, // 属性下发结果上报
	// 	"getStatusReply": _ithings_PropertyUpTopic, // 获取实时状态上报
	// 	"eventPost":      _ithings_EventUpTopic,    // 事件上报
	// 子设备上报
	if msg.Method == "subDevReport" {
		if !tc.topologyReady {
			return 0, fmt.Errorf("sub device topology not ready")
		}
		msg.Method = "report"
		if subDev, ok := tc.subDevices[msg.SubDeviceName]; ok {
			msg.Timestamp = time.Now().UnixMilli()                      // 上报时间戳使用毫秒
			msg.ClientToken = fmt.Sprintf("%d", time.Now().UnixMicro()) // Token 使用时微秒间戳
			err := tc.client.Publish(subDev.ReportTopic(), 1, false, msg.string()).Error()
			if err != nil {
				glogger.GLogger.Error(err)
			}
		}
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
type TopologyDevice struct {
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

// $gateway/status/${productid}/${devicename}
func (tt TopologyDevice) OnlineTopic() string {
	return fmt.Sprintf("$gateway/status/%s/%s", tt.ProductID, tt.DeviceName)
}

// $gateway/status/${productid}/${devicename}
func (tt TopologyDevice) OfflineTopic() string {
	return fmt.Sprintf("$gateway/status/%s/%s", tt.ProductID, tt.DeviceName)
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
