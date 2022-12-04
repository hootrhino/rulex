package target

import (
	"net/http"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

type HTTPTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig common.HTTPConfig
	status     typex.SourceState
}

func NewHTTPTarget(e typex.RuleX) typex.XTarget {
	ht := new(HTTPTarget)
	ht.RuleEngine = e
	ht.mainConfig = common.HTTPConfig{}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *HTTPTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}

	return nil

}
func (ht *HTTPTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
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
	return ht.status

}
func (ht *HTTPTarget) To(data interface{}) (interface{}, error) {
	r, err := utils.Post(ht.client, data, ht.mainConfig.Url, ht.mainConfig.Headers)
	return r, err
}

func (ht *HTTPTarget) Stop() {
	ht.status = typex.SOURCE_STOP
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
	return core.GenOutConfig(typex.HTTP_TARGET, "HTTP_TARGET", common.HTTPConfig{})
}
