// uart_driver相当于是升级版，这个是最原始的基础驱动
package driver

import (
	"context"
	"errors"
	"sync"

	"github.com/hootrhino/rulex/typex"
	serial "github.com/wwhai/goserial"
)

type rawUartDriver struct {
	state      typex.DriverState
	serialPort serial.Port
	ctx        context.Context
	RuleEngine typex.RuleX
	device     *typex.Device
	lock       sync.Mutex
}

// 初始化一个驱动
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
		lock:       sync.Mutex{},
	}
}

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
	a.lock.Lock()
	_, err := a.serialPort.Write([]byte("\r\n"))
	a.lock.Unlock()
	return err
}

func (a *rawUartDriver) Read(cmd []byte, b []byte) (int, error) {
	a.lock.Lock()
	n, e := a.serialPort.Read(b)
	a.lock.Unlock()
	return n, e
}

func (a *rawUartDriver) Write(cmd []byte, b []byte) (int, error) {
	return a.serialPort.Write(b)
}
func (a *rawUartDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Raw Uart Driver",
		Type:        "RAW_UART",
		Description: "Raw Uart Driver",
	}
}
