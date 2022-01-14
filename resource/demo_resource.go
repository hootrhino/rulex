package resource

import "rulex/typex"

type DemoResource struct {
	typex.XStatus
}

func (*DemoResource) Driver() typex.XExternalDriver {
	return nil
}
func (m *DemoResource) OnStreamApproached(data string) error {
	return nil
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

func (d *DemoResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
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
func (*DemoResource) Configs() typex.XConfig {
	return typex.XConfig{}
}

//
// 拓扑
//
func (*DemoResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
