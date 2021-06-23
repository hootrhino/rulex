package x

import (
	"fmt"
	"time"

	"github.com/ngaut/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
type MqttInEndResource struct {
	inEndId string
	client  mqtt.Client
}

func NewMqttInEndResource(inEndId string) *MqttInEndResource {
	return &MqttInEndResource{
		inEndId: inEndId,
	}
}

func (mm *MqttInEndResource) Start(e *RuleEngine, successCallBack func(), errorCallback func(error)) error {

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Infof("[RULE Engine Log] ===> Received message: [%s] from topic: [%s]\n", msg.Payload(), msg.Topic())
		e.Work(GetInEnd(mm.inEndId), string(msg.Payload()))
	}
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("[RULE Engine Log] ===> Mqtt InEnd Connected Success")
		// TODO support multipul topics
		client.Subscribe("$X_IN_END", 1, nil)
		// Update Running Resource State
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Infof("[RULE Engine Log] ===> Connect lost: %v\n", err)
		// Update Running Resource State
	}
	log.Debug(GetInEnd(mm.inEndId))
	config := GetInEnd(mm.inEndId).Config
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", (*config)["server"], (*config)["port"].(int)))
	opts.SetClientID("x-client-main-" + (*config)["clientId"].(string))
	opts.SetUsername((*config)["username"].(string))
	opts.SetPassword((*config)["password"].(string))
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(5 * time.Second)
	mm.client = mqtt.NewClient(opts)
	if token := mm.client.Connect(); token.Wait() && token.Error() != nil {
		errorCallback(token.Error())
		return token.Error()
	} else {
		successCallBack()
		return nil
	}

}
func (mm *MqttInEndResource) Stop() {

}
func (mm *MqttInEndResource) Reload() {

}
func (mm *MqttInEndResource) Pause() {

}
func (mm *MqttInEndResource) Status() int {
	return GetInEnd(mm.inEndId).State
}

func (mm *MqttInEndResource) Register(inEndId string) error {

	return nil
}
