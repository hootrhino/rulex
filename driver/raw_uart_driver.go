//
//
// uart_driver相当于是升级版，这个是最原始的基础驱动
//
//
package driver

import (
	"context"
	"errors"

	"github.com/i4de/rulex/typex"
	serial "github.com/wwhai/goserial"
)

type rawUartDriver struct {
	state      typex.DriverState
	config     serial.Config
	serialPort serial.Port
	ctx        context.Context
	RuleEngine typex.RuleX
	device     *typex.Device
}

//
// 初始化一个驱动
//
func NewRawUartDriver(
	ctx context.Context,
	e typex.RuleX,
	device *typex.Device,
	serialPort serial.Port,
) typex.XExternalDriver {
	return &rawUartDriver{
		RuleEngine: e,
		ctx:        ctx,
		serialPort: serialPort,
		device:     device,
	}
}

//
//
//
func (a *rawUartDriver) Init(map[string]string) error {
	a.state = typex.DRIVER_UP

	return nil
}

func (a *rawUartDriver) Work() error {

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
func (a *rawUartDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Raw Uart Driver",
		Type:        "RAW_UART",
		Description: "Raw Uart Driver",
	}
}
