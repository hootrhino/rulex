package core

import (
	"github.com/gin-gonic/gin"
)

//
type HttpInEndResource struct {
	*XStatus
	engine *gin.Engine
	e      *RuleEngine
}

func NewHttpInEndResource(inEndId string, e *RuleEngine) *HttpInEndResource {

	h := HttpInEndResource{}
	h.InEndId = inEndId
	h.engine = gin.Default()
	h.e = e
	return &h
}

//
func (hh *HttpInEndResource) Start() error {
	hh.engine = gin.New()
	config := hh.e.GetInEnd(hh.InEndId).Config
	hh.engine.GET("/in", func(c *gin.Context) {
		inForm := struct{ data string }{}
		err := c.BindJSON(inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err,
			})
		} else {
			hh.e.Work(hh.e.GetInEnd(hh.InEndId), inForm.data)
			c.JSON(200, gin.H{
				"message": "ok",
				"data":    inForm,
			})
		}
	})
	hh.engine.Run(":" + (*config)["port"].(string))
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
	return hh.e.GetInEnd(hh.InEndId).State
}

func (hh *HttpInEndResource) Register(inEndId string) error {

	return nil
}

func (hh *HttpInEndResource) Test(inEndId string) bool {
	return true
}

func (hh *HttpInEndResource) Enabled() bool {
	return hh.Enable
}
