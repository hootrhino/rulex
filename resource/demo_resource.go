package resource

import "rulex/typex"

type DemoResource struct {
	typex.XStatus
}

func (d *DemoResource) Details() *typex.InEnd {
	return d.RuleEngine.GetInEnd(d.PointId)
}

func (d *DemoResource) Test(inEndId string) bool {
	return true
}

func (d *DemoResource) Register(inEndId string) error {
	return nil
}

func (d *DemoResource) Start() error {
	return nil

}

func (d *DemoResource) Enabled() bool {
	return true
}

func (d *DemoResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
}

func (d *DemoResource) Reload() {
}

func (d *DemoResource) Pause() {
}

func (d *DemoResource) Status() typex.ResourceState {
	return typex.UP
}

func (d *DemoResource) Stop() {
}
