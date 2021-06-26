package x

import (
	"github.com/gin-gonic/gin"
)

//
type HttpInEndResource struct {
	enabled bool
	inEndId string
	engine  *gin.Engine
}

func NewHttpInEndResource(inEndId string) *HttpInEndResource {
	return &HttpInEndResource{
		inEndId: inEndId,
		engine:  gin.Default(),
	}
}
func (hh *HttpInEndResource) Start(e *RuleEngine) error {
	hh.engine = gin.New()
	gin.SetMode(gin.ReleaseMode)
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
func (hh *HttpInEndResource) Status(e *RuleEngine) int {
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
