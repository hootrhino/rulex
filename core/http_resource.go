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
	engine     *gin.Engine
	ruleEngine *RuleEngine
}

func NewHttpInEndResource(inEndId string, e *RuleEngine) *HttpInEndResource {
	h := HttpInEndResource{}
	h.InEndId = inEndId
	h.engine = gin.Default()
	h.ruleEngine = e
	return &h
}

//
func (hh *HttpInEndResource) Start() error {
	config := hh.ruleEngine.GetInEnd(hh.InEndId).Config
	hh.engine.GET("/in", func(c *gin.Context) {
		inForm := struct{ data string }{}
		err := c.BindJSON(inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err,
			})
		} else {
			hh.ruleEngine.Work(hh.ruleEngine.GetInEnd(hh.InEndId), inForm.data)
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
	return hh.ruleEngine.GetInEnd(hh.InEndId).State
}

func (hh *HttpInEndResource) Register(inEndId string) error {
	hh.InEndId = inEndId
	return nil
}

func (hh *HttpInEndResource) Test(inEndId string) bool {
	return true
}

func (hh *HttpInEndResource) Enabled() bool {
	return hh.Enable
}
