package httpserver

import (
	"net/http"

	"github.com/i4de/rulex/typex"

	"github.com/gin-gonic/gin"
)

const SUCCESS string = "Success"

// Http Return
type R struct {
	Code int         `json:"code" binding:"required"`
	Msg  string      `json:"msg" binding:"required"`
	Data interface{} `json:"data"`
}

func Ok() R {
	return R{200, SUCCESS, []interface{}{}}
}
func OkWithEmpty() R {
	return R{200, SUCCESS, []interface{}{}}
}
func OkWithData(data interface{}) R {
	return R{200, SUCCESS, data}
}
func Error(s string) R {
	return R{4000, s, nil}
}
func Error400(e error) R {
	return R{4001, e.Error(), []interface{}{}}
}
func Error500(e error) R {
	return R{5001, e.Error(), []interface{}{}}
}

func (hh *HttpApiServer) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

//
//
//
func Cros() gin.HandlerFunc {
	return func(c *gin.Context) {
		cros(c)
	}
}

//
func cros(c *gin.Context) {
	c.Header("Cache-Control", "private, max-age=10")
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	c.Header("Access-Control-Max-Age", "172800")
	c.Header("Access-Control-Allow-Credentials", "true")

	if method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Request.Header.Del("Origin")
	c.Next()
}

//
// Add api route
//
func (h *HttpApiServer) addRoute(f func(*gin.Context, *HttpApiServer, typex.RuleX)) func(*gin.Context) {

	return func(c *gin.Context) {
		f(c, h, h.ruleEngine)
	}
}
