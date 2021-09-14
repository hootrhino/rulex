package core

type DemoTarget struct {
	XStatus
}

func (d *DemoTarget) Details() *outEnd {
	return d.RuleEngine.GetOutEnd(d.PointId)
}

func (d *DemoTarget) Test(outEndId string) bool {
	return true
}

func (d *DemoTarget) Register(outEndId string) error {
	return nil
}

func (d *DemoTarget) Start() error {
	return nil

}

func (d *DemoTarget) Enabled() bool {
	return true
}

func (d *DemoTarget) Reload() {
}

func (d *DemoTarget) Pause() {
}

func (d *DemoTarget) Status() State {
	return UP
}

func (d *DemoTarget) To(data interface{}) error {
	return nil

}

func (d *DemoTarget) Stop() {
}
