package driver

import (
	"rulex/typex"

	"github.com/goburrow/modbus"
)

/*
*
* Modbus RTU 驱动直接用了库，所以这个驱动仅仅是为了符合模式，其实没有实际作用，或者留着以后扩展用
*
 */
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
		state:      typex.RUNNING,
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

func (d *modBusRtuDriver) Read(data []byte) (int, error) {
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
