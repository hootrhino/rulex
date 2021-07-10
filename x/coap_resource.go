package x

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
}

func NewCoAPInEndResource(inEndId string) *CoAPInEndResource {
	c := CoAPInEndResource{}
	c.InEndId = inEndId
	c.router = mux.NewRouter()
	return &c
}

func (cc *CoAPInEndResource) Start(e *RuleEngine) error {

	cc.router.Use(func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			log.Debugf("ClientAddress %v, %v\n", w.Client().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	})
	cc.router.Handle("/in", mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		log.Debugf("Received Coap Data: %#v", msg)
		e.Work(e.GetInEnd(cc.InEndId), msg.String())
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

func (cc *CoAPInEndResource) Reload() {

}
func (cc *CoAPInEndResource) Pause() {

}
func (cc *CoAPInEndResource) Status(e *RuleEngine) State {
	return UP
}

func (cc *CoAPInEndResource) Register(inEndId string) error {
	cc.InEndId = inEndId
	return nil
}

func (cc *CoAPInEndResource) Test(inEndId string) bool {
	return true
}
func (cc *CoAPInEndResource) Enabled() bool {
	return cc.Enable
}
