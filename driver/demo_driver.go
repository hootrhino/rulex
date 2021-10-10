package driver

import "rulex/typex"

type DemoDriver struct {
}

func (d *DemoDriver) Test() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Init() error {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Work() error {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) State() typex.DriverState {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Stop() error {
	panic("not implemented") // TODO: Implement
}
