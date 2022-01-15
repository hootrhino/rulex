package resource

import (
	"context"
	"fmt"
	"net/http"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

//
type httpConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}

//
type httpInEndResource struct {
	typex.XStatus
	engine *gin.Engine
}

func NewHttpInEndResource(inEndId string, e typex.RuleX) typex.XResource {
	h := httpInEndResource{}
	h.PointId = inEndId
	gin.SetMode(gin.ReleaseMode)
	h.engine = gin.New()
	h.RuleEngine = e
	return &h
}
func (*httpInEndResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("HTTP", "", httpConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

//
func (hh *httpInEndResource) Start() error {
	config := hh.RuleEngine.GetInEnd(hh.PointId).Config
	var mainConfig httpConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
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
			hh.RuleEngine.Work(hh.RuleEngine.GetInEnd(hh.PointId), inForm.Data)
			c.JSON(200, gin.H{
				"message": "ok",
				"data":    inForm,
			})
		}
	})
	go func(ctx context.Context) {
		http.ListenAndServe(fmt.Sprintf(":%v", mainConfig.Port), hh.engine)
	}(context.Background())
	log.Info("HTTP resource started on" + " [0.0.0.0]:" + fmt.Sprintf("%v", mainConfig.Port))

	return nil
}

//
func (mm *httpInEndResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

//
func (hh *httpInEndResource) Stop() {

}
func (hh *httpInEndResource) Reload() {

}
func (hh *httpInEndResource) Pause() {

}
func (hh *httpInEndResource) Status() typex.ResourceState {
	return typex.UP
}

func (hh *httpInEndResource) Register(inEndId string) error {
	hh.PointId = inEndId
	return nil
}

func (hh *httpInEndResource) Test(inEndId string) bool {
	return true
}

func (hh *httpInEndResource) Enabled() bool {
	return hh.Enable
}
func (hh *httpInEndResource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}
func (m *httpInEndResource) OnStreamApproached(data string) error {
	return nil
}
func (*httpInEndResource) Driver() typex.XExternalDriver {
	return nil
}

//
// 拓扑
//
func (*httpInEndResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
