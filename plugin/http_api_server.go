package plugin

import (
	"context"
	"net/http"
	"rulenginex/statistics"
	"rulenginex/x"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const API_ROOT string = "/api/v1/"
const DASHBOARD_ROOT string = "/dashboard/v1/"

type HttpApiServer struct {
	ginEngine  *gin.Engine
	RuleEngine *x.RuleEngine
}

func (hh *HttpApiServer) Load(r *x.RuleEngine) *x.XPluginEnv {
	hh.ginEngine = gin.New()
	gin.SetMode(gin.ReleaseMode)
	hh.ginEngine.LoadHTMLGlob("plugin/templates/*")
	hh.RuleEngine = r
	return x.NewXPluginEnv()
}

//
func (hh *HttpApiServer) Init(env *x.XPluginEnv) error {

	ctx := context.Background()
	go func(ctx context.Context) {
		hh.ginEngine.Run(":2580")
	}(ctx)
	return nil
}
func (hh *HttpApiServer) Install(env *x.XPluginEnv) (*x.XPluginMetaInfo, error) {
	return &x.XPluginMetaInfo{
		Name:     "HttpApiServer",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}, nil
}

//
//
func (hh *HttpApiServer) Start(e *x.RuleEngine, env *x.XPluginEnv) error {
	hh.ginEngine.GET(DASHBOARD_ROOT, func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{})
	})
	hh.ginEngine.GET(API_ROOT+"plugins", func(c *gin.Context) {
		cros(c)
		c.PureJSON(http.StatusOK, gin.H{
			"plugins": hh.RuleEngine.Plugins,
		})
	})
	hh.ginEngine.GET(API_ROOT+"system", func(c *gin.Context) {
		cros(c)
		//
		percent, _ := cpu.Percent(time.Second, false)
		memInfo, _ := mem.VirtualMemory()
		parts, _ := disk.Partitions(true)
		diskInfo, _ := disk.Usage(parts[0].Mountpoint)
		c.JSON(http.StatusOK, gin.H{
			"diskInfo":   diskInfo.UsedPercent,
			"memInfo":    memInfo.UsedPercent,
			"cpuPercent": percent[0],
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"cpus":       runtime.GOMAXPROCS(0)})
	})
	//
	hh.ginEngine.GET(API_ROOT+"inends", func(c *gin.Context) {
		cros(c)
		c.JSON(http.StatusOK, gin.H{"inends": e.AllInEnd()})
	})
	//
	hh.ginEngine.GET(API_ROOT+"outends", func(c *gin.Context) {
		cros(c)
		c.JSON(http.StatusOK, gin.H{"outends": e.AllOutEnd()})
	})
	//
	hh.ginEngine.GET(API_ROOT+"rules", func(c *gin.Context) {
		cros(c)
		c.JSON(http.StatusOK, gin.H{"rules": x.AllRule()})
	})
	//
	hh.ginEngine.GET(API_ROOT+"statistics", func(c *gin.Context) {
		cros(c)
		c.JSON(http.StatusOK, gin.H{"statistics": statistics.AllStatistics()})
	})
	//
	log.Info("Http web dashboard started on:http://127.0.0.1:2580" + DASHBOARD_ROOT)
	return nil
}

func (hh *HttpApiServer) Uninstall(env *x.XPluginEnv) error {
	log.Info("HttpApiServer Uninstalled")
	return nil
}
func (hh *HttpApiServer) Clean() {
	log.Info("HttpApiServer Cleaned")
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
