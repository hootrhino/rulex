package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

type myClient struct{}

func Test_104_client(t *testing.T) {
	var err error

	option := cs104.NewOption()
	if err = option.AddRemoteServer("127.0.0.1:2404"); err != nil {
		panic(err)
	}

	mycli := &myClient{}

	client := cs104.NewClient(mycli, option)

	client.LogMode(true)

	client.SetOnConnectHandler(func(c *cs104.Client) {
		c.SendStartDt() // 发送startDt激活指令
	})
	err = client.Start()
	if err != nil {
		panic(fmt.Errorf("Failed to connect. error:%v\n", err))
	}

	for {
		time.Sleep(time.Second * 100)
	}

}
func (myClient) InterrogationHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (myClient) ReadHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) TestCommandHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) ClockSyncHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (myClient) ResetProcessHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (myClient) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (myClient) ASDUHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
