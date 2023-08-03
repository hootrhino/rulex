package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CorsMiddleWare 跨域中间件
// CORS是跨源资源分享（Cross-Origin Resource Sharing）中间件
func CorsMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//指定允许其他域名访问
		//ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*") //跨域：CORS(跨来源资源共享)策略
		//预检结果缓存时间
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		//允许的请求类型（GET,POST等）
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		//允许的请求头字段
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		//是否允许后续请求携带认证信息（cookies）,该值只能是true,否则不返回
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}
