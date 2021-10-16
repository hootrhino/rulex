package resource

import (
	"bytes"
	"context"
	"fmt"
	"rulex/typex"

	"github.com/ngaut/log"
	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

//
type CoAPInEndResource struct {
	typex.XStatus
	router *mux.Router
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

	var port = ""
	switch (*config)["port"].(type) {
	case string:
		port = ":" + (*config)["port"].(string)
		break
	case int:
		port = fmt.Sprintf(":%v", (*config)["port"].(int))
		break
	case int64:
		port = fmt.Sprintf(":%v", (*config)["port"].(int))
		break
	case float64:
		port = fmt.Sprintf(":%v", (*config)["port"].(int))
		break
	}
	cc.router.Use(func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			log.Debugf("Client Address %v, %v\n", w.Client().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	})
	cc.router.Handle("/in", mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		log.Debugf("Received Coap Data: %#v", msg)
		cc.RuleEngine.Work(cc.RuleEngine.GetInEnd(cc.PointId), msg.String())
		err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("ok")))
		if err != nil {
			log.Errorf("Cannot set response: %v", err)
		}
	}))
	go func(ctx context.Context) {
		err := coap.ListenAndServe("udp", port, cc.router)
		if err != nil {
			return
		} else {
			return
		}
	}(context.Background())
	cc.Enable = true
	log.Info("Coap resource started on [udp]" + port)
	return nil
}
func (m *CoAPInEndResource) OnStreamApproached(data string) error {
	return nil
}
//
func (cc *CoAPInEndResource) Stop() {
}

func (mm *CoAPInEndResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
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
