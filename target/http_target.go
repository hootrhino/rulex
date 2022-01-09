package target

import (
	"net/http"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
)

type httpConfig struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
type HTTPTarget struct {
	typex.XStatus
	url     string
	headers map[string]string
	client  http.Client
}

func NewHTTPTarget(e typex.RuleX) typex.XTarget {
	ht := new(HTTPTarget)
	ht.RuleEngine = e
	return ht
}
func (ht *HTTPTarget) Register(outEndId string) error {
	ht.PointId = outEndId
	return nil

}
func (ht *HTTPTarget) Start() error {
	config := ht.RuleEngine.GetOutEnd(ht.PointId).Config
	var mainConfig httpConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	ht.url = mainConfig.Url
	ht.headers = mainConfig.Headers
	ht.client = http.Client{}
	log.Info("HTTPTarget started")
	return nil
}
func (ht *HTTPTarget) OnStreamApproached(data string) error {
	_, err := utils.Post(ht.client, data, ht.url, ht.headers)
	return err
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
	_, err := utils.Post(ht.client, data, ht.url, ht.headers)
	return err
}

//
func (ht *HTTPTarget) Stop() {
}
func (ht *HTTPTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}
