package source

import (
	"bytes"
	"context"
	"fmt"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"

	"github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/mux"
)

//
type coAPConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}

//
type coAPInEndSource struct {
	typex.XStatus
	router     *mux.Router
	port       uint16
	dataModels []typex.XDataModel
}

func NewCoAPInEndSource(inEndId string, e typex.RuleX) *coAPInEndSource {
	c := coAPInEndSource{}
	c.PointId = inEndId
	c.router = mux.NewRouter()
	c.RuleEngine = e
	return &c
}

func (cc *coAPInEndSource) Start(cctx typex.CCTX) error {
	cc.Ctx = cctx.Ctx
	cc.CancelCTX = cctx.CancelCTX

	config := cc.RuleEngine.GetInEnd(cc.PointId).Config
	var mainConfig coAPConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	port := fmt.Sprintf(":%v", mainConfig.Port)
	cc.dataModels = mainConfig.DataModels
	cc.router.Use(func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			// glogger.GLogger.Debugf("Client Address %v, %v\n", w.Client().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	})
	//
	// /in
	//
	cc.router.Handle("/in", mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		// glogger.GLogger.Debugf("Received Coap Data: %#v", msg)
		work, err := cc.RuleEngine.WorkInEnd(cc.RuleEngine.GetInEnd(cc.PointId), msg.String())
		if !work {
			glogger.GLogger.Error(err)
		}
		if err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("ok"))); err != nil {
			glogger.GLogger.Errorf("Cannot set response: %v", err)
		}
	}))

	go func(ctx context.Context) {
		err := coap.ListenAndServe("udp", port, cc.router)
		if err != nil {
			glogger.GLogger.Error(err)
			return
		}
	}(cc.Ctx)
	glogger.GLogger.Info("Coap source started on [udp]" + port)
	return nil
}

//
func (cc *coAPInEndSource) Stop() {
	cc.CancelCTX()
}

func (cc *coAPInEndSource) DataModels() []typex.XDataModel {
	return cc.XDataModels
}
func (cc *coAPInEndSource) Reload() {

}
func (cc *coAPInEndSource) Pause() {

}
func (cc *coAPInEndSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (cc *coAPInEndSource) Init(inEndId string, cfg map[string]interface{}) error {
	cc.PointId = inEndId
	var mainConfig coAPConfig
	if err := utils.BindSourceConfig(cfg, &mainConfig); err != nil {
		return err
	}
	cc.port = mainConfig.Port
	cc.dataModels = mainConfig.DataModels
	return nil
}
func (cc *coAPInEndSource) Test(inEndId string) bool {
	return true
}
func (cc *coAPInEndSource) Enabled() bool {
	return true
}
func (cc *coAPInEndSource) Details() *typex.InEnd {
	return cc.RuleEngine.GetInEnd(cc.PointId)
}

func (cc *coAPInEndSource) Driver() typex.XExternalDriver {
	return nil
}

func (*coAPInEndSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.COAP, "COAP", coAPConfig{})
}

//
// 拓扑
//
func (*coAPInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
