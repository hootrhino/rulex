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

// 即将废弃
func (s *RulexApiServer) AddRoute(f func(c *gin.Context,
	ruleEngine typex.RuleX)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c, s.ruleEngine)
	}
}

// New api after 0.6.4
func AddRoute(f func(c *gin.Context,
	ruleEngine typex.RuleX)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c, DefaultApiServer.ruleEngine)
	}
}

// AddRouteV2: Only for cron，It's redundant API， will refactor in feature
func AddRouteV2(f func(*gin.Context, typex.RuleX) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		data, err := f(c, DefaultApiServer.ruleEngine)
		if err != nil {
			c.JSON(response.HTTP_OK, response.Error400(err))
		} else {
			c.JSON(response.HTTP_OK, response.OkWithData(data))
		}
	}
}

func (s *RulexApiServer) GetGroup(name string) *gin.RouterGroup {
	return s.ginEngine.Group(name)
}

// new API
func RouteGroup(name string) *gin.RouterGroup {
	return DefaultApiServer.ginEngine.Group(name)
}
func (s *RulexApiServer) Route() *gin.Engine {
	return s.ginEngine
}

/*
*
* 初始化网络配置, 后期分一个windows版本
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
	// 初始化默认路由, 如果没有配置会在数据库生成关于eth1的一个默认路由数据
	service.InitDefaultIpRoute()
	// 一组操作, 主要用来初始化 DHCP和DNS、网卡配置等
	// 1 2 3 的目的是为了每次重启的时候初始化软路由
	{
		// 1 初始化默认路由表: ip route
		service.ConfigDefaultIpTable()
		// 2 初始化默认DHCP
		service.ConfigDefaultDhcp()
		// 3 初始化Eth1的静态IP地址
		service.ConfigDefaultNat()
	}
	// 配置一个默认分组
	service.InitGenericGroup(&model.MGenericGroup{
		UUID:   "ROOT",
		Type:   "VISUAL",
		Name:   "默认分组",
		Parent: "NULL",
	})
	service.InitGenericGroup(&model.MGenericGroup{
		UUID:   "ROOT",
		Type:   "DEVICE",
		Name:   "默认分组",
		Parent: "NULL",
	})
}
