package test

import (
	"github.com/eclipse/paho.mqtt.golang"

	"testing"
)

func publish(t *testing.T) {

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		t.Log("Mqtt InEnd Connected Success")
		client.Publish("$x_IN_END", 1, false, "test is ok1")
		client.Publish("$x_IN_END", 1, false, "test is ok2")
		client.Publish("$x_IN_END", 1, false, "test is ok3")
		t.Log("Publish to x_IN_END ok")

	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		t.Log("Connect lost:", err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	//
	opts.SetClientID("x_IN_END_TEST1")
	opts.SetUsername("x_IN_END_TEST1")
	opts.SetPassword("x_IN_END_TEST1")
	//
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Log("error:", token.Error())
	}
}
