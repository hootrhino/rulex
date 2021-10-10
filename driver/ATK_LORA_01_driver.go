package driver

import (
	"context"
	"rulex/typex"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 串口数据写入
//
func write(a *ATK_LORA_01Driver, k string) (string, error) {
	_, err := a.serialPort.Write([]byte(k + "\r\n"))
	if err != nil {
		return "", err
	}
	for {
		response := make([]byte, 4)
		size, err := a.serialPort.Read(response)
		if err != nil {
			return "", err
		}
		if size > 0 {
			return string(response), nil
		}
	}
}

//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type ATK_LORA_01Driver struct {
	serialPort *serial.Port
	channel    chan bool
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

//
// 初始化一个驱动
//
func NewATK_LORA_01Driver(serialPort *serial.Port, in *typex.InEnd, e typex.RuleX) typex.XDriver {
	m := &ATK_LORA_01Driver{}
	m.channel = make(chan bool)
	m.In = in
	m.RuleEngine = e
	m.serialPort = serialPort
	m.ctx = context.Background()
	return m
}

//
//
//
func (a *ATK_LORA_01Driver) Init() error {
	return nil
}
func (a *ATK_LORA_01Driver) Work() error {
	go func(context.Context) {
		log.Debug("ATK LORA 01 Driver Start Listening...")
		for {
			select {
			case <-a.ctx.Done():
				return
			default:
				{
					response := make([]byte, 16) // byte
					size, err := a.serialPort.Read(response)
					if err != nil {
						a.Stop()
						return
					} else {
						// log.Debug("SerialPort Received:", string(response))
						err := a.RuleEngine.PushQueue(typex.QueueData{
							In:   a.In,
							Out:  nil,
							E:    a.RuleEngine,
							Data: string(response[:size]),
						})
						if err != nil {
							log.Error("ATK_LORA_01Driver error: ", err)
						}
					}
				}
			}

		}
	}(a.ctx)
	return nil

}
func (a *ATK_LORA_01Driver) State() typex.DriverState {
	return typex.RUNNING

}
func (a *ATK_LORA_01Driver) Stop() error {
	a.ctx.Done()
	return nil
}

func (a *ATK_LORA_01Driver) Test() (string, error) {
	return write(a, "AT\r\n")
}
