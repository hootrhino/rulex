package source

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"

	"github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/mux"
)

type coAPInEndSource struct {
	typex.XStatus
	router     *mux.Router
	mainConfig common.HostConfig
	status     typex.SourceState
}

func NewCoAPInEndSource(e typex.RuleX) typex.XSource {
	c := coAPInEndSource{}
	c.router = mux.NewRouter()
	c.mainConfig = common.HostConfig{}
	c.RuleEngine = e
	return &c
}

func (cc *coAPInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	cc.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &cc.mainConfig); err != nil {
		return err
	}

	return nil
}
func (cc *coAPInEndSource) Start(cctx typex.CCTX) error {
	cc.Ctx = cctx.Ctx
	cc.CancelCTX = cctx.CancelCTX
	port := fmt.Sprintf(":%v", cc.mainConfig.Port)
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
	cc.status = typex.SOURCE_UP
	glogger.GLogger.Info("Coap source started on [udp]" + port)
	return nil
}

func (cc *coAPInEndSource) Stop() {
	cc.status = typex.SOURCE_STOP
	if cc.CancelCTX != nil {
		cc.CancelCTX()
	}
}

func (cc *coAPInEndSource) DataModels() []typex.XDataModel {
	return cc.XDataModels
}

func (cc *coAPInEndSource) Status() typex.SourceState {
	return cc.status
}

func (cc *coAPInEndSource) Test(inEndId string) bool {
	return true
}

func (cc *coAPInEndSource) Details() *typex.InEnd {
	return cc.RuleEngine.GetInEnd(cc.PointId)
}

func (cc *coAPInEndSource) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*coAPInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*coAPInEndSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*coAPInEndSource) UpStream([]byte) (int, error) {
	return 0, nil
}
