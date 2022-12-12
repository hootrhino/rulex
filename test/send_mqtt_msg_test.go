package test

import (
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Test_publish_msg_to_rulex(t *testing.T) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("127.0.0.1:1883")
	opts.SetClientID("go-unit-test-1")
	opts.SetUsername("go-unit-test-1")
	opts.SetPassword("go-unit-test-1")

	opts.SetPingTimeout(3 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.WaitTimeout(3 * time.Second)
	if token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	} else {
		//
		//
		// 发送一个打开开关的指令下去:
		// 这个指令是硬件支持的协议格式: {"cmdId": "00001","cmd" :"on","sw": [1, 2] }
		//
		client.Publish("rulex-client-topic-1", 2, false, `{"cmdId": "00001","cmd" :"on","sw": [1, 2] }`)
		time.Sleep(time.Second)
		client.Publish("rulex-client-topic-1", 2, false, `{"cmdId": "00001","cmd" :"off","sw": [1, 2] }`)
	}

	time.Sleep(time.Second)
}
