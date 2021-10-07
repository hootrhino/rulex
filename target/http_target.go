package target

import (
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
)

type HTTPTarget struct {
	typex.XStatus
	url string
}

func (ht *HTTPTarget) Register(outEndId string) {
	ht.PointId = outEndId
	config := ht.RuleEngine.GetOutEnd(ht.PointId).Config
	ht.url = (*config)["url"].(string)

}
func (ht *HTTPTarget) Start() {
	ht.Enable = true
	log.Info("HTTPTarget started")
}
func (ht *HTTPTarget) Test(outEndId string) bool {
	return true
}
func (ht *HTTPTarget) Enabled() bool {
	return ht.Enable
}
func (ht *HTTPTarget) Reload() {

}
func (ht *HTTPTarget) Pause() {

}
func (ht *HTTPTarget) Status() typex.ResourceState {
	return typex.UP

}
func (ht *HTTPTarget) To(data interface{}) error {
	_, err := utils.Post(data, ht.url)
	return err
}

//
func (ht *HTTPTarget) Stop() {

}
