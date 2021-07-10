package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Query("token") // 访问令牌

		if token != "" {
			// TODO add jwt Authorize support
			c.Next()
		} else {
			// c.Abort()
			// c.JSON(http.StatusUnauthorized, gin.H{"msg": "No authority operate"})
			c.Next()
			return
		}
	}
}

//
func cros(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
	}
}
