package target

import (
	"errors"
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
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
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
	if mm.client != nil {
		mm.client.Disconnect(0)
	}
	mm.CancelCTX()
	mm.status = typex.SOURCE_DOWN

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

//
//
//
func (mm *mqttOutEndTarget) To(data interface{}) (interface{}, error) {
	if mm.client != nil {
		return mm.client.Publish(mm.mainConfig.PubTopic, 1, false, data).Error(), nil
	}
	return nil, errors.New("mqtt client is nil")
}

/*
*
* 配置
*
 */
func (*mqttOutEndTarget) Configs() *typex.XConfig {
	return core.GenOutConfig(typex.MQTT_TARGET, "MQTT", common.MqttConfig{})
}
