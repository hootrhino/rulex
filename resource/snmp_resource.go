package resource

import "rulex/typex"

type SNMPResource struct {
	typex.XStatus
}

func NewSNMPInEndResource(inEndId string, e typex.RuleX) *SNMPResource {
	u := SNMPResource{}

	return &u
}
func (s *SNMPResource) Details() *typex.InEnd {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Test(inEndId string) bool {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Register(inEndId string) error {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Start() error {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Enabled() bool {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) DataModels() *map[string]typex.XDataModel {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Reload() {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Pause() {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Status() typex.ResourceState {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) OnStreamApproached(data string) error {
	panic("not implemented") // TODO: Implement
}

func (s *SNMPResource) Stop() {
	panic("not implemented") // TODO: Implement
}
