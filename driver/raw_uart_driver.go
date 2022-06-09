//
//
// uart_driver相当于是升级版，这个是最原始的基础驱动
//
//
package driver

import (
	"context"
	"errors"
	"rulex/typex"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

type rawUartDriver struct {
	state      typex.DriverState
	serialPort serial.Port
	config     serial.Config
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

//
// 初始化一个驱动
//
func NewRawUartDriver(
	ctx context.Context,
	config serial.Config,
	in *typex.InEnd,
	e typex.RuleX,
	onRead func([]byte)) (typex.XExternalDriver, error) {

	return &rawUartDriver{
		In:         in,
		RuleEngine: e,
		config:     config,
		ctx:        ctx,
	}, nil
}

//
//
//
func (a *rawUartDriver) Init(map[string]string) error {
	a.state = typex.DRIVER_RUNNING
	return nil
}

func (a *rawUartDriver) Work() error {
	serialPort, err := serial.Open(&a.config)
	a.serialPort = serialPort
	if err != nil {
		log.Error("uartModuleSource start failed:", err)
		return err
	}
	return nil

}
func (a *rawUartDriver) State() typex.DriverState {
	return a.state
}
func (a *rawUartDriver) Stop() error {
	a.state = typex.DRIVER_STOP
	return a.serialPort.Close()
}

func (a *rawUartDriver) Test() error {
	if a.serialPort == nil {
		return errors.New("serialPort is nil")
	}
	_, err := a.serialPort.Write([]byte("\r\n"))
	return err

}

//
func (a *rawUartDriver) Read(b []byte) (int, error) {
	return a.serialPort.Read(b)
}

//
func (a *rawUartDriver) Write(b []byte) (int, error) {
	return a.serialPort.Write(b)
}
func (a *rawUartDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "Raw Uart Driver",
		Type:        "RAW_UART",
		Description: "Raw Uart Driver",
	}
}
