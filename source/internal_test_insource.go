package source

import "github.com/hootrhino/rulex/typex"

type internalTestSource struct {
	typex.XStatus
}

func NewInternalTestSource(e typex.RuleX) typex.XSource {
	return &internalTestSource{}
}
func (*internalTestSource) Configs() *typex.XConfig {
	return &typex.XConfig{}
}
func (hh *internalTestSource) Init(inEndId string, configMap map[string]interface{}) error {
	hh.PointId = inEndId
	return nil
}

func (hh *internalTestSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	return nil
}

func (mm *internalTestSource) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

func (hh *internalTestSource) Stop() {
	hh.CancelCTX()
}
func (hh *internalTestSource) Reload() {

}
func (hh *internalTestSource) Pause() {

}
func (hh *internalTestSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (hh *internalTestSource) Test(inEndId string) bool {
	return true
}

func (hh *internalTestSource) Enabled() bool {
	return hh.Enable
}
func (hh *internalTestSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}

func (*internalTestSource) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*internalTestSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*internalTestSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*internalTestSource) UpStream([]byte) (int, error) {
	return 0, nil
}
