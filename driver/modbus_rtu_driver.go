package driver

import (
	"rulex/typex"

	"github.com/goburrow/modbus"
)

type modBusRtuDriver struct {
	state      typex.DriverState
	client     modbus.Client
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

func NewModBusRtuDriver(
	in *typex.InEnd,
	e typex.RuleX,
	client modbus.Client) typex.XExternalDriver {
	return &modBusRtuDriver{
		In:         in,
		RuleEngine: e,
		client:     client,
	}

}
func (d *modBusRtuDriver) Test() error {
	return nil
}

func (d *modBusRtuDriver) Init() error {
	return nil
}

func (d *modBusRtuDriver) Work() error {
	return nil
}

func (d *modBusRtuDriver) State() typex.DriverState {
	return d.state
}

func (d *modBusRtuDriver) SetState(s typex.DriverState) {
	d.state = s
}

func (d *modBusRtuDriver) Read(_ []byte) (int, error) {
	return 0, nil

}

func (d *modBusRtuDriver) Write(_ []byte) (int, error) {
	return 0, nil
}

func (d *modBusRtuDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "ModBus RTU Driver",
		Type:        "UART",
		Description: "ModBus RTU Driver",
	}
}

func (d *modBusRtuDriver) Stop() error {
	return nil
}
