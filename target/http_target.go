package target

import (
	"net/http"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
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
func (ht *HTTPTarget) Init(outEndId string, cfg map[string]interface{}) error {
	ht.PointId = outEndId
	return nil

}
func (ht *HTTPTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	config := ht.RuleEngine.GetOutEnd(ht.PointId).Config
	var mainConfig httpConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	ht.url = mainConfig.Url
	ht.headers = mainConfig.Headers
	ht.client = http.Client{}
	glogger.GLogger.Info("HTTPTarget started")
	return nil
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
func (ht *HTTPTarget) Status() typex.SourceState {
	return typex.SOURCE_UP

}
func (ht *HTTPTarget) To(data interface{}) (interface{}, error) {
	r, err := utils.Post(ht.client, data, ht.url, ht.headers)
	return r, err
}

//
func (ht *HTTPTarget) Stop() {
	ht.CancelCTX()
}
func (ht *HTTPTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}

/*
*
* 配置
*
 */
func (*HTTPTarget) Configs() *typex.XConfig {
	return core.GenOutConfig(typex.HTTP_TARGET, "HTTP_TARGET", httpConfig{})
}
