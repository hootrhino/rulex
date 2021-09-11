package core

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

//
type HttpInEndResource struct {
	XStatus
	engine *gin.Engine
}

func NewHttpInEndResource(inEndId string, e RuleX) *HttpInEndResource {
	h := HttpInEndResource{}
	h.PointId = inEndId
	h.engine = gin.Default()
	h.RuleEngine = e
	return &h
}

//
func (hh *HttpInEndResource) Start() error {
	config := hh.RuleEngine.GetInEnd(hh.PointId).Config
	hh.engine.GET("/in", func(c *gin.Context) {
		inForm := struct{ data string }{}
		err := c.BindJSON(inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			hh.RuleEngine.Work(hh.RuleEngine.GetInEnd(hh.PointId), inForm.data)
			c.JSON(200, gin.H{
				"message": "ok",
				"data":    inForm,
			})
		}
	})
	go func(ctx context.Context) {
		http.ListenAndServe(":"+(*config)["port"].(string), hh.engine)
	}(context.Background())
	log.Info("HTTP Resource start successfully")

	return nil
}

//
func (mm *HttpInEndResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{}
}

//
func (hh *HttpInEndResource) Stop() {

}
func (hh *HttpInEndResource) Reload() {

}
func (hh *HttpInEndResource) Pause() {

}
func (hh *HttpInEndResource) Status() State {
	return hh.RuleEngine.GetInEnd(hh.PointId).State
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
func (hh *HttpInEndResource) Details() *inEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}
