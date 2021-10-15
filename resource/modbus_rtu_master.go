package resource

import "rulex/typex"

type ModbusRTUMasterResource struct {
	typex.XStatus
}

func NewModbusRTUMasterResource(inEndId string, e typex.RuleX) typex.XResource {
	s := ModbusRTUMasterResource{}
	s.PointId = inEndId
	s.RuleEngine = e
	//
	return &s
}

func (m *ModbusRTUMasterResource) Register(inEndId string) error {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Start() error {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Details() *typex.InEnd {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Test(inEndId string) bool {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Enabled() bool {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) DataModels() *map[string]typex.XDataModel {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Reload() {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Pause() {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Status() typex.ResourceState {
	panic("not implemented") // TODO: Implement
}

func (m *ModbusRTUMasterResource) Stop() {
	panic("not implemented") // TODO: Implement
}
