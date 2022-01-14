package resource

import (
	"bytes"
	"context"
	"fmt"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
	"github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

//
type coAPConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}

//
type CoAPInEndResource struct {
	typex.XStatus
	router     *mux.Router
	dataModels []typex.XDataModel
}

func NewCoAPInEndResource(inEndId string, e typex.RuleX) *CoAPInEndResource {
	c := CoAPInEndResource{}
	c.PointId = inEndId
	c.router = mux.NewRouter()
	c.RuleEngine = e
	return &c
}

func (cc *CoAPInEndResource) Start() error {
	config := cc.RuleEngine.GetInEnd(cc.PointId).Config
	var mainConfig coAPConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	port := fmt.Sprintf(":%v", mainConfig.Port)
	cc.dataModels = mainConfig.DataModels
	cc.router.Use(func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			// log.Debugf("Client Address %v, %v\n", w.Client().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	})
	//
	// /in
	//
	cc.router.Handle("/in", mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		// log.Debugf("Received Coap Data: %#v", msg)
		cc.RuleEngine.Work(cc.RuleEngine.GetInEnd(cc.PointId), msg.String())
		err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("ok")))
		if err != nil {
			log.Errorf("Cannot set response: %v", err)
		}
	}))
	go func(ctx context.Context) {
		err := coap.ListenAndServe("udp", port, cc.router)
		if err != nil {
			log.Error(err)
			return
		}
	}(context.Background())
	log.Info("Coap resource started on [udp]" + port)
	return nil
}
func (m *CoAPInEndResource) OnStreamApproached(data string) error {
	return nil
}

//
func (cc *CoAPInEndResource) Stop() {
}

func (cc *CoAPInEndResource) DataModels() []typex.XDataModel {
	return cc.dataModels
}
func (cc *CoAPInEndResource) Reload() {

}
func (cc *CoAPInEndResource) Pause() {

}
func (cc *CoAPInEndResource) Status() typex.ResourceState {
	return typex.UP
}

func (cc *CoAPInEndResource) Register(inEndId string) error {
	cc.PointId = inEndId
	return nil
}

func (cc *CoAPInEndResource) Test(inEndId string) bool {
	return true
}
func (cc *CoAPInEndResource) Enabled() bool {
	return true
}
func (cc *CoAPInEndResource) Details() *typex.InEnd {
	return cc.RuleEngine.GetInEnd(cc.PointId)
}

func (cc *CoAPInEndResource) Driver() typex.XExternalDriver {
	return nil
}

func (*CoAPInEndResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("COAP", "", coAPConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

//
// 拓扑
//
func (*CoAPInEndResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
