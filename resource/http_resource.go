package resource

import (
	"context"
	"net/http"
	"rulex/typex"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

//
type HttpInEndResource struct {
	typex.XStatus
	engine *gin.Engine
}

func NewHttpInEndResource(inEndId string, e typex.RuleX) *HttpInEndResource {
	h := HttpInEndResource{}
	h.PointId = inEndId
	h.engine = gin.Default()
	h.RuleEngine = e
	return &h
}

//
func (hh *HttpInEndResource) Start() error {
	config := hh.RuleEngine.GetInEnd(hh.PointId).Config
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
		http.ListenAndServe(":"+(*config)["port"].(string), hh.engine)
	}(context.Background())
	log.Info("HTTP resource started on" + " [0.0.0.0]:" + (*config)["port"].(string))

	return nil
}

//
func (mm *HttpInEndResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
}

//
func (hh *HttpInEndResource) Stop() {

}
func (hh *HttpInEndResource) Reload() {

}
func (hh *HttpInEndResource) Pause() {

}
func (hh *HttpInEndResource) Status() typex.ResourceState {
	return typex.UP
}

func (hh *HttpInEndResource) Register(inEndId string) error {
	hh.PointId = inEndId
	return nil
}

func (hh *HttpInEndResource) Test(inEndId string) bool {
	return true
}

func (hh *HttpInEndResource) Enabled() bool {
	return hh.Enable
}
func (hh *HttpInEndResource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}
func (m *HttpInEndResource) OnStreamApproached(data string) error {
	return nil
}