package source

import (
	"context"
	"fmt"
	"net/http"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/gin-gonic/gin"
)

//
type httpConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}

//
type httpInEndSource struct {
	typex.XStatus
	engine *gin.Engine
}

func NewHttpInEndSource(inEndId string, e typex.RuleX) typex.XSource {
	h := httpInEndSource{}
	h.PointId = inEndId
	gin.SetMode(gin.ReleaseMode)
	h.engine = gin.New()
	h.RuleEngine = e
	return &h
}
func (*httpInEndSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.HTTP, "HTTP", httpConfig{})
}

//
func (hh *httpInEndSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	config := hh.RuleEngine.GetInEnd(hh.PointId).Config
	var mainConfig httpConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	hh.XDataModels = mainConfig.DataModels
	hh.engine.POST("/in", func(c *gin.Context) {
		type Form struct {
			Data string
		}
		var inForm Form
		err := c.BindJSON(&inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			hh.RuleEngine.WorkInEnd(hh.RuleEngine.GetInEnd(hh.PointId), inForm.Data)
			c.JSON(200, gin.H{
				"message": "ok",
				"data":    inForm,
			})
		}
	})

	go func(ctx context.Context) {
		err := http.ListenAndServe(fmt.Sprintf(":%v", mainConfig.Port), hh.engine)
		if err != nil {
			glogger.GLogger.Error(err)
			return
		}
	}(hh.Ctx)
	glogger.GLogger.Info("HTTP source started on" + " [0.0.0.0]:" + fmt.Sprintf("%v", mainConfig.Port))

	return nil
}

//
func (mm *httpInEndSource) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

//
func (hh *httpInEndSource) Stop() {
	hh.CancelCTX()

}
func (hh *httpInEndSource) Reload() {

}
func (hh *httpInEndSource) Pause() {

}
func (hh *httpInEndSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (hh *httpInEndSource) Init(inEndId string, cfg map[string]interface{}) error {
	hh.PointId = inEndId
	return nil
}
func (hh *httpInEndSource) Test(inEndId string) bool {
	return true
}

func (hh *httpInEndSource) Enabled() bool {
	return hh.Enable
}
func (hh *httpInEndSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}

func (*httpInEndSource) Driver() typex.XExternalDriver {
	return nil
}

//
// 拓扑
//
func (*httpInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
