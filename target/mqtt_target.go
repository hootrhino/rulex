package target

import (
	"errors"
	"fmt"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//

type mqttOutEndTarget struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig common.MqttConfig
	status     typex.SourceState
}

func NewMqttTarget(e typex.RuleX) typex.XTarget {
	m := new(mqttOutEndTarget)
	m.RuleEngine = e
	m.mainConfig = common.MqttConfig{}
	m.status = typex.SOURCE_DOWN
	return m
}
func (*mqttOutEndTarget) Driver() typex.XExternalDriver {
	return nil
}
func (mm *mqttOutEndTarget) Init(outEndId string, configMap map[string]interface{}) error {
	mm.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &mm.mainConfig); err != nil {
		return err
	}
	return nil
}
func (mm *mqttOutEndTarget) Start(cctx typex.CCTX) error {
	mm.Ctx = cctx.Ctx
	mm.CancelCTX = cctx.CancelCTX
	//
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("Mqtt OutEnd Connected Success")
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
	opts.SetAutoReconnect(false)
	opts.SetMaxReconnectInterval(0)
	mm.client = mqtt.NewClient(opts)
	token := mm.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		mm.status = typex.SOURCE_UP
		return nil
	}

}

func (mm *mqttOutEndTarget) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

func (mm *mqttOutEndTarget) Stop() {
	mm.CancelCTX()
	mm.status = typex.SOURCE_DOWN
	if mm.client != nil {
		mm.client.Disconnect(0)
		mm.client = nil
	}
}
func (mm *mqttOutEndTarget) Reload() {

}
func (mm *mqttOutEndTarget) Pause() {

}
func (mm *mqttOutEndTarget) Status() typex.SourceState {
	return mm.status
}

func (mm *mqttOutEndTarget) Test(outEndId string) bool {
	if mm.client != nil {
		return mm.client.IsConnected()
	}
	return false
}

func (mm *mqttOutEndTarget) Enabled() bool {
	return mm.Enable
}
func (mm *mqttOutEndTarget) Details() *typex.OutEnd {
	return mm.RuleEngine.GetOutEnd(mm.PointId)
}

func (mm *mqttOutEndTarget) To(data interface{}) (interface{}, error) {
	if mm.client != nil {
		token := mm.client.Publish(mm.mainConfig.PubTopic, 1, false, data)
		return token.Error(), nil
	}
	return nil, errors.New("mqtt client is nil")
}

/*
*
* 配置
*
 */
func (*mqttOutEndTarget) Configs() *typex.XConfig {
	return &typex.XConfig{}
}
