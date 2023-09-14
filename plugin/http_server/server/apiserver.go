package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/glogger"
	response "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 拼接URL
*
 */
func ContextUrl(path string) string {
	return API_V1_ROOT + path
}

const API_V1_ROOT string = "/api/v1/"
const DEFAULT_DB_PATH string = "./rulex.db"

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

/*
*
* 开启Server
*
 */
func StartRulexApiServer(ruleEngine typex.RuleX) {
	gin.SetMode(gin.ReleaseMode)
	server := RulexApiServer{
		ginEngine:  gin.New(),
		ruleEngine: ruleEngine,
		config:     serverConfig{Port: 2580},
	}
	server.ginEngine.Use(static.Serve("/", WWWRoot("")))
	server.ginEngine.Use(Authorize())
	server.ginEngine.Use(Cros())
	server.ginEngine.GET("/ws", glogger.WsLogger)
	server.ginEngine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		glogger.GLogger.Error(err)
		c.JSON(200, response.Error500(err1crash))
	}))
	/*
	*
	* 解决浏览器刷新被重定向问题
	*
	 */
	server.ginEngine.NoRoute(func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.JSON(404, response.Error("No such Route:"+c.Request.URL.Path))
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Write(indexHTML)
		c.Writer.Flush()
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
func (s *RulexApiServer) AddRoute(f func(c *gin.Context,
	ruleEngine typex.RuleX)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c, s.ruleEngine)
	}
}

func (s *RulexApiServer) GetGroup(name string) *gin.RouterGroup {
	return s.ginEngine.Group(name)
}
func (s *RulexApiServer) Route() *gin.Engine {
	return s.ginEngine
}

/*
*
* 初始化网络配置
*
 */
func (s *RulexApiServer) InitializeData() {
	// 加载资源类型
	source.LoadSt()
	target.LoadTt()
	device.LoadDt()
	// 初始化有线网口配置
	if !service.CheckIfAlreadyInitNetWorkConfig() {
		service.InitNetWorkConfig()
	}
	// 初始化WIFI配置
	if !service.CheckIfAlreadyInitWlanConfig() {
		service.InitWlanConfig()
	}
	// 初始化网站配置
	service.InitSiteConfig(model.MSiteConfig{
		SiteName: "RhinoEEKIT",
		Logo:     "RhinoEEKIT",
		AppName:  "RhinoEEKIT",
	})
	// 初始化默认路由
	service.InitDefaultIpRoute()
}
