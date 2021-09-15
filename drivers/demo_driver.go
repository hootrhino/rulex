package drivers

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

func (d *DemoDriver) State() DriverState {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Stop() error {
	panic("not implemented") // TODO: Implement
}
