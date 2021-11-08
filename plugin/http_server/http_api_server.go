package httpserver

import (
	"context"
	"rulex/typex"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"

	"gorm.io/gorm"
)

const API_ROOT string = "/api/v1/"
const DASHBOARD_ROOT string = "/dashboard/v1/"

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
func (hh *HttpApiServer) Load() *typex.XPluginEnv {
	return typex.NewXPluginEnv()
}

//
func (hh *HttpApiServer) Init() error {
	gin.SetMode(gin.ReleaseMode)
	hh.ginEngine = gin.New()
	hh.ginEngine.Use(Authorize())
	hh.ginEngine.Use(Cros())
	if hh.dbPath == "" {
		hh.InitDb("./rulex.db")
	} else {
		hh.InitDb(hh.dbPath)
	}
	ctx := context.Background()
	go func(ctx context.Context, port int) {
		hh.ginEngine.Run(":" + strconv.Itoa(port))
	}(ctx, hh.Port)
	return nil
}

//
// HttpApiServer Start
//
func (hh *HttpApiServer) Start() error {

	//
	// Get all plugins
	//
	hh.ginEngine.GET(API_ROOT+"plugins", hh.addRoute(Plugins))
	//
	// Get system infomation
	//
	hh.ginEngine.GET(API_ROOT+"system", hh.addRoute(System))
	//
	// Get all inends
	//
	hh.ginEngine.GET(API_ROOT+"inends", hh.addRoute(InEnds))
	//
	// Get all outends
	//
	hh.ginEngine.GET(API_ROOT+"outends", hh.addRoute(OutEnds))
	//
	// Get all rules
	//
	hh.ginEngine.GET(API_ROOT+"rules", hh.addRoute(Rules))
	//
	// Get statistics data
	//
	hh.ginEngine.GET(API_ROOT+"statistics", hh.addRoute(Statistics))
	//
	// Auth
	//
	hh.ginEngine.POST(API_ROOT+"users", hh.addRoute(CreateUser))
	hh.ginEngine.POST(API_ROOT+"auth", hh.addRoute(Auth))
	//
	// Create InEnd
	//
	hh.ginEngine.POST(API_ROOT+"inends", hh.addRoute(CreateInend))
	//
	// Create OutEnd
	//
	hh.ginEngine.POST(API_ROOT+"outends", hh.addRoute(CreateOutEnd))
	//
	// Create rule
	//
	hh.ginEngine.POST(API_ROOT+"rules", hh.addRoute(CreateRule))
	//
	// Delete inend by UUID
	//
	hh.ginEngine.DELETE(API_ROOT+"inends", hh.addRoute(DeleteInend))
	//
	// Delete outend by UUID
	//
	hh.ginEngine.DELETE(API_ROOT+"outends", hh.addRoute(DeleteOutend))
	//
	// Delete rule by UUID
	//
	hh.ginEngine.DELETE(API_ROOT+"rules", hh.addRoute(DeleteRule))
	//
	log.Info("Http server started on http://127.0.0.1:2580")
	return nil
}

func (hh *HttpApiServer) Stop() error {
	return nil
}


func (hh *HttpApiServer) Db() *gorm.DB {
	return hh.sqliteDb
}
func (hh *HttpApiServer) XPluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "HttpApiServer",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
