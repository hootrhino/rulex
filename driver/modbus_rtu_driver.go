package driver

import (
	"rulex/typex"

	"github.com/goburrow/modbus"
)

type ModBusRtuDriver struct {
	state      typex.DriverState
	client     modbus.Client
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

func NewModBusRtuDriver(
	in *typex.InEnd,
	e typex.RuleX,
	client modbus.Client) typex.XExternalDriver {
	return &ModBusRtuDriver{
		In:         in,
		RuleEngine: e,
		client:     client,
	}

}
func (d *ModBusRtuDriver) Test() error {
	return nil
}

func (d *ModBusRtuDriver) Init() error {
	return nil
}

func (d *ModBusRtuDriver) Work() error {
	return nil
}

func (d *ModBusRtuDriver) State() typex.DriverState {
	return d.state
}

func (d *ModBusRtuDriver) SetState(s typex.DriverState) {
	d.state = s
}

func (d *ModBusRtuDriver) Read(_ []byte) (int, error) {
	return 0, nil

}

func (d *ModBusRtuDriver) Write(_ []byte) (int, error) {
	return 0, nil
}

func (d *ModBusRtuDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "ModBus RTU Driver",
		Type:        "UART",
		Description: "ModBus RTU Driver",
	}
}

func (d *ModBusRtuDriver) Stop() error {
	return nil
}
