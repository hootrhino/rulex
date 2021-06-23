package x

import (
	"bytes"

	"github.com/ngaut/log"
	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

//
type CoAPInEndResource struct {
	inEndId string
	router  *mux.Router
}

func NewCoAPInEndResource(inEndId string) *CoAPInEndResource {
	return &CoAPInEndResource{
		inEndId: inEndId,
		router:  mux.NewRouter(),
	}
}

func (cc *CoAPInEndResource) Start(e *RuleEngine, successCallBack func(), errorCallback func(error)) error {

	cc.router.Use(func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			log.Debugf("ClientAddress %v, %v\n", w.Client().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	})
	cc.router.Handle("/in", mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		log.Debugf("Received Coap Data: %#v", msg)
		e.Work(GetInEnd(cc.inEndId), msg.String())
		err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("ok")))
		if err != nil {
			log.Errorf("cannot set response: %v", err)
		}
	}))
	err := coap.ListenAndServe("udp", ":5688", cc.router)
	if err != nil {
		errorCallback(err)
		return err
	} else {
		successCallBack()
		return nil
	}
}

//
func (cc *CoAPInEndResource) Stop() {

}

func (cc *CoAPInEndResource) Reload() {

}
func (cc *CoAPInEndResource) Pause() {

}
func (cc *CoAPInEndResource) Status() int {
	return GetInEnd(cc.inEndId).State
}

func (cc *CoAPInEndResource) Register(inEndId string) error {

	return nil
}
