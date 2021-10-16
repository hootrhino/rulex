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
func NewATK_LORA_01Driver(serialPort *serial.Port, in *typex.InEnd, e typex.RuleX) typex.XExternalDriver {
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
				}
			}

			response1 := make([]byte, 1)   // byte
			response2 := make([]byte, 128) // byte
			size1, err1 := a.serialPort.Read(response1)
			size2, err2 := a.serialPort.Read(response2)

			if err1 != nil || err2 != nil {
				err := a.Stop()
				if err != nil {
					return
				}
				log.Error("ATK_LORA_01Driver error: ", err1, err2)
				return
			} else {
				response := string(append(response1, response2...))
				//log.Debug("SerialPort Received:", size1+size2)
				err0 := a.RuleEngine.PushQueue(typex.QueueData{
					In:   a.In,
					Out:  nil,
					E:    a.RuleEngine,
					Data: response[:(size1 + size2)],
				})
				if err0 != nil {
					log.Error("ATK_LORA_01Driver error: ", err0)
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

func (a *ATK_LORA_01Driver) Test() error {
	return nil
}

//
func (a *ATK_LORA_01Driver) Read([]byte) (int, error) {

	return 0, nil
}

//
func (a *ATK_LORA_01Driver) Write([]byte) (int, error) {
	return 0, nil
}
