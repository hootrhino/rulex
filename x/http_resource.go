package x

import (
	"github.com/gin-gonic/gin"
)

//
type HttpInEndResource struct {
	*XStatus
	engine *gin.Engine
}

func NewHttpInEndResource(inEndId string) *HttpInEndResource {
	h := HttpInEndResource{}
	h.inEndId = inEndId
	h.engine = gin.Default()
	return &h
}
func (hh *HttpInEndResource) Start(e *RuleEngine) error {
	hh.engine = gin.New()
	config := e.GetInEnd(hh.inEndId).Config
	hh.engine.GET("/in", func(c *gin.Context) {
		inForm := struct{ data string }{}
		err := c.BindJSON(inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err,
			})
		} else {
			e.Work(e.GetInEnd(hh.inEndId), inForm.data)
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
func (hh *HttpInEndResource) Stop() {

}
func (hh *HttpInEndResource) Reload() {

}
func (hh *HttpInEndResource) Pause() {

}
func (hh *HttpInEndResource) Status(e *RuleEngine) State {
	return e.GetInEnd(hh.inEndId).State
}

func (hh *HttpInEndResource) Register(inEndId string) error {

	return nil
}

func (hh *HttpInEndResource) Test(inEndId string) bool {
	return true
}

func (hh *HttpInEndResource) Enabled() bool {
	return hh.enabled
}
