package driver

import (
	"context"
	"errors"

	"github.com/i4de/rulex/typex"
	serial "github.com/wwhai/goserial"
)

type UsrG776Driver struct {
	state      typex.DriverState
	serialPort serial.Port
	ctx        context.Context
	RuleEngine typex.RuleX
	device     *typex.Device
}

// 初始化一个驱动
func NewUsrG776Driver(
	ctx context.Context,
	e typex.RuleX,
	device *typex.Device,
	serialPort serial.Port,
) typex.XExternalDriver {
	return &UsrG776Driver{
		RuleEngine: e,
		ctx:        ctx,
		serialPort: serialPort,
		device:     device,
	}
}

func (d *UsrG776Driver) Init(map[string]string) error {
	d.state = typex.DRIVER_UP

	return nil
}

func (d *UsrG776Driver) Work() error {

	return nil

}
func (d *UsrG776Driver) State() typex.DriverState {
	return d.state
}
func (d *UsrG776Driver) Stop() error {
	d.state = typex.DRIVER_STOP
	return d.serialPort.Close()
}

func (d *UsrG776Driver) Test() error {
	if d.serialPort == nil {
		return errors.New("serialPort is nil")
	}
	_, err := d.serialPort.Write([]byte("\r\n"))
	return err

}

func (d *UsrG776Driver) Read(cmd int, b []byte) (int, error) {

	return d.serialPort.Read(b)
}

func (d *UsrG776Driver) Write(cmd int, b []byte) (int, error) {
	return d.serialPort.Write(b)
}
func (d *UsrG776Driver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "UsrG776Driver Driver",
		Type:        "RAW_UART",
		Description: "UsrG776Driver Driver",
	}
}
