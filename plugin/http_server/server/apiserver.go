package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/hootrhino/rulex/glogger"
	response "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

// 全局API Server
var DefaultApiServer *RulexApiServer

/*
*
* API Server
*
 */
type RulexApiServer struct {
	ginEngine  *gin.Engine
	ruleEngine typex.RuleX
	config     serverConfig
}
type serverConfig struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

var err1crash = errors.New("http server crash, try to recovery")

func StartRulexApiServer(ruleEngine typex.RuleX) {
	gin.SetMode(gin.ReleaseMode)
	server := RulexApiServer{
		ginEngine:  gin.New(),
		ruleEngine: ruleEngine,
		config:     serverConfig{Port: 2580},
	}
	server.ginEngine.Use(Authorize())
	server.ginEngine.Use(Cros())
	server.ginEngine.Use(static.Serve("/", WWWRoot("")))
	server.ginEngine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		glogger.GLogger.Error(err)
		c.JSON(200, response.Error500(err1crash))
	}))
	server.ginEngine.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})
	//
	// Http server
	//
	go func(ctx context.Context, port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
		if err := server.ginEngine.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
	}(typex.GCTX, server.config.Port)
	DefaultApiServer = &server
}
func (s *RulexApiServer) AddRoute(f func(c *gin.Context, ruleEngine typex.RuleX)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c, s.ruleEngine)
	}
}
