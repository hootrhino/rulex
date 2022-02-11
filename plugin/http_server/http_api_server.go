package httpserver

import (
	"context"
	"net/http"
	"rulex/typex"
	"rulex/utils"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"

	"gorm.io/gorm"
)

const _API_V1_ROOT string = "/api/v1/"
const DEFAULT_DB_PATH string = "./rulex.db"
const DASHBOARD_ROOT string = "/dashboard/v1/"

type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
}
type HttpApiServer struct {
	Port       int
	Root       string
	sqliteDb   *gorm.DB
	dbPath     string
	ginEngine  *gin.Engine
	ruleEngine typex.RuleX
}

func NewHttpApiServer(port int, root string, dbPath string, e typex.RuleX) *HttpApiServer {
	return &HttpApiServer{Port: port, Root: root, dbPath: dbPath, ruleEngine: e}
}

//
func (hh *HttpApiServer) Init(cfg interface{}) error {
	gin.SetMode(gin.ReleaseMode)
	hh.ginEngine = gin.New()
	configHttpServer(hh)
	go func(ctx context.Context, port int) {
		hh.ginEngine.Run(":" + strconv.Itoa(port))
	}(typex.GCTX, hh.Port)
	return nil
}

//
// HttpApiServer Start
//
func (hh *HttpApiServer) Start() error {
	hh.ginEngine.GET("/", hh.addRoute(Index))

	//
	// Get all plugins
	//
	hh.ginEngine.GET(_API_V1_ROOT+"plugins", hh.addRoute(Plugins))
	//
	// Get system infomation
	//
	hh.ginEngine.GET(_API_V1_ROOT+"system", hh.addRoute(System))
	//
	// Ping -> Pong
	//
	hh.ginEngine.GET(_API_V1_ROOT+"ping", hh.addRoute(Ping))
	//
	//
	//
	hh.ginEngine.GET(_API_V1_ROOT+"sourceCount", hh.addRoute(SourceCount))
	//
	//
	//
	hh.ginEngine.GET(_API_V1_ROOT+"logs", hh.addRoute(Logs))
	//
	//
	//
	hh.ginEngine.POST(_API_V1_ROOT+"logout", hh.addRoute(LogOut))
	//
	// Get all inends
	//
	hh.ginEngine.GET(_API_V1_ROOT+"inends", hh.addRoute(InEnds))
	//
	//
	//
	hh.ginEngine.GET(_API_V1_ROOT+"drivers", hh.addRoute(Drivers))
	//
	// Get all outends
	//
	hh.ginEngine.GET(_API_V1_ROOT+"outends", hh.addRoute(OutEnds))
	//
	// Get all rules
	//
	hh.ginEngine.GET(_API_V1_ROOT+"rules", hh.addRoute(Rules))
	//
	// Get statistics data
	//
	hh.ginEngine.GET(_API_V1_ROOT+"statistics", hh.addRoute(Statistics))
	//
	// Auth
	//
	hh.ginEngine.POST(_API_V1_ROOT+"users", hh.addRoute(CreateUser))
	//
	//
	//
	hh.ginEngine.POST(_API_V1_ROOT+"login", hh.addRoute(Login))
	//
	//
	//
	hh.ginEngine.GET(_API_V1_ROOT+"info", hh.addRoute(Info))
	//
	// Create InEnd
	//
	hh.ginEngine.POST(_API_V1_ROOT+"inends", hh.addRoute(CreateInend))
	hh.ginEngine.GET(_API_V1_ROOT+"inends/config", hh.addRoute(GetInEndConfig))
	//
	// Create OutEnd
	//
	hh.ginEngine.POST(_API_V1_ROOT+"outends", hh.addRoute(CreateOutEnd))
	//
	// Create rule
	//
	hh.ginEngine.POST(_API_V1_ROOT+"rules", hh.addRoute(CreateRule))
	//
	// Delete inend by UUID
	//
	hh.ginEngine.DELETE(_API_V1_ROOT+"inends", hh.addRoute(DeleteInend))
	//
	// Delete outend by UUID
	//
	hh.ginEngine.DELETE(_API_V1_ROOT+"outends", hh.addRoute(DeleteOutend))
	//
	// Delete rule by UUID
	//
	hh.ginEngine.DELETE(_API_V1_ROOT+"rules", hh.addRoute(DeleteRule))
	//
	// 验证 lua 语法
	//
	hh.ginEngine.POST(_API_V1_ROOT+"validateRule", hh.addRoute(ValidateLuaSyntax))
	//
	// 获取配置表
	//
	hh.ginEngine.GET(_API_V1_ROOT+"rType", hh.addRoute(RType))
	hh.ginEngine.GET(_API_V1_ROOT+"tType", hh.addRoute(TType))
	//
	// 串口列表
	//
	hh.ginEngine.GET(_API_V1_ROOT+"getUartList", hh.addRoute(GetUartList))

	log.Info("Http server started on http://127.0.0.1:", hh.Port)
	return nil
}

func (hh *HttpApiServer) Stop() error {
	return nil
}

func (hh *HttpApiServer) Db() *gorm.DB {
	return hh.sqliteDb
}
func (hh *HttpApiServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "Http Api Server",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
func configHttpServer(hh *HttpApiServer) {
	hh.ginEngine.Use(Authorize())
	hh.ginEngine.Use(Cros())
	hh.ginEngine.LoadHTMLFiles(utils.GetPwd() + hh.Root + "index.html")
	hh.ginEngine.StaticFS("/static", http.Dir(utils.GetPwd()+hh.Root+"static"))
	hh.ginEngine.StaticFS("/assets", http.Dir(utils.GetPwd()+hh.Root+"static/assets"))
	hh.ginEngine.StaticFile("/favicon.ico", utils.GetPwd()+hh.Root+"favicon.ico")
	if hh.dbPath == "" {
		hh.InitDb(DEFAULT_DB_PATH)
	} else {
		hh.InitDb(hh.dbPath)
	}
}
