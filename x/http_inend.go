package x

import (
	"github.com/gin-gonic/gin"
)

//
type HttpInEndResource struct {
	inEndId string
	engine  *gin.Engine
}

func NewHttpInEndResource(inEndId string) *HttpInEndResource {
	return &HttpInEndResource{
		inEndId: inEndId,
		engine:  gin.Default(),
	}
}
func (hh *HttpInEndResource) Start(e *RuleEngine, successCallBack func(), errorCallback func(error)) error {
	hh.engine = gin.Default()
	config := GetInEnd(hh.inEndId).Config

	hh.engine.GET("/in", func(c *gin.Context) {
		inForm := struct{ data string }{}
		err := c.BindJSON(inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err,
			})
		} else {
			e.Work(GetInEnd(hh.inEndId), inForm.data)
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
func (hh *HttpInEndResource) Status() int {
	return GetInEnd(hh.inEndId).State
}

func (hh *HttpInEndResource) Register(inEndId string) error {

	return nil
}
