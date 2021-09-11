package core

import (
	"bytes"
	"context"

	"github.com/ngaut/log"
	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

//
type CoAPInEndResource struct {
	XStatus
	router *mux.Router
	e      *RuleEngine
}

func NewCoAPInEndResource(inEndId string, e *RuleEngine) *CoAPInEndResource {
	c := CoAPInEndResource{}
	c.PointId = inEndId
	c.router = mux.NewRouter()
	c.e = e
	return &c
}

func (cc *CoAPInEndResource) Start() error {

	cc.router.Use(func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			log.Debugf("ClientAddress %v, %v\n", w.Client().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	})
	cc.router.Handle("/in", mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		log.Debugf("Received Coap Data: %#v", msg)
		cc.e.Work(cc.e.GetInEnd(cc.PointId), msg.String())
		err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("ok")))
		if err != nil {
			log.Errorf("cannot set response: %v", err)
		}
	}))
	go func(ctx context.Context) {
		err := coap.ListenAndServe("udp", ":5688", cc.router)
		if err != nil {
			return
		} else {
			return
		}
	}(context.Background())
	cc.Enable = true

	return nil
}

//
func (cc *CoAPInEndResource) Stop() {

}

func (mm *CoAPInEndResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{}
}
func (cc *CoAPInEndResource) Reload() {

}
func (cc *CoAPInEndResource) Pause() {

}
func (cc *CoAPInEndResource) Status() State {
	return UP
}

func (cc *CoAPInEndResource) Register(inEndId string) error {
	cc.PointId = inEndId
	return nil
}

func (cc *CoAPInEndResource) Test(inEndId string) bool {
	return true
}
func (cc *CoAPInEndResource) Enabled() bool {
	return cc.Enable
}
func (cc *CoAPInEndResource) Details() *inEnd {
	return cc.RuleEngine.GetInEnd(cc.PointId)
}
