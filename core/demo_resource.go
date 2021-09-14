package core

type DemoResource struct {
	XStatus
}

func (d *DemoResource) Details() *inEnd {
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

func (d *DemoResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{}
}

func (d *DemoResource) Reload() {
}

func (d *DemoResource) Pause() {
}

func (d *DemoResource) Status() ResourceState {
	return UP
}

func (d *DemoResource) Stop() {
}
