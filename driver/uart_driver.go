package driver

import (
	"context"
	"rulex/typex"
	"time"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type UartDriver struct {
	serialPort *serial.Port
	channel    chan bool
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

//
// 初始化一个驱动
//
func NewUartDriver(serialPort *serial.Port, in *typex.InEnd, e typex.RuleX) typex.XExternalDriver {
	m := &UartDriver{}
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
func (a *UartDriver) Init() error {
	return nil
}
func (a *UartDriver) Work() error {
	go func(context.Context) {
		ticker := time.NewTicker(30 * time.Microsecond)
		log.Debug("UartDriver Start Listening")
		for {
			select {
			case <-a.ctx.Done():
				return
			default:
				{
					<-ticker.C
					response1 := make([]byte, 1)   // byte
					response2 := make([]byte, 128) // byte
					size1, err1 := a.serialPort.Read(response1)
					size2, err2 := a.serialPort.Read(response2)
					if err1 != nil || err2 != nil {
						err := a.Stop()
						if err != nil {
							return
						}
						log.Error("UartDriver error: ", err1, err2)
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
							log.Error("UartDriver error: ", err0)
						}
					}
				}
			}

		}
	}(a.ctx)
	return nil

}
func (a *UartDriver) State() typex.DriverState {
	return typex.RUNNING

}
func (a *UartDriver) Stop() error {
	a.ctx.Done()
	return nil
}

func (a *UartDriver) Test() error {
	return nil
}

//
func (a *UartDriver) Read([]byte) (int, error) {

	return 0, nil
}

//
func (a *UartDriver) Write([]byte) (int, error) {
	return 0, nil
}
