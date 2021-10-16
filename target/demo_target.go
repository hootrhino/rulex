package target

import "rulex/typex"

type DemoTarget struct {
	typex.XStatus
}
func (m *DemoTarget) OnStreamApproached(data string) error {
	return nil
}
func (d *DemoTarget) Details() *typex.OutEnd {
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

func (d *DemoTarget) Status() typex.ResourceState {
	return typex.UP
}

func (d *DemoTarget) To(data interface{}) error {
	return nil

}

func (d *DemoTarget) Stop() {
}
