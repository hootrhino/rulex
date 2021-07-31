package httpserver

import (
	"context"
	"rulex/core"
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
	ginEngine  *gin.Engine
	ruleEngine *core.RuleEngine
}

func NewHttpApiServer(port int, root string, e *core.RuleEngine) *HttpApiServer {
	return &HttpApiServer{Port: port, Root: root, ruleEngine: e}
}
func (hh *HttpApiServer) Load() *core.XPluginEnv {
	return core.NewXPluginEnv()
}

//
func (hh *HttpApiServer) Init(env *core.XPluginEnv) error {
	gin.SetMode(gin.ReleaseMode)
	hh.ginEngine = gin.New()
	hh.ginEngine.Use(Authorize())
	hh.InitDb()
	hh.ginEngine.LoadHTMLGlob(hh.Root)
	ctx := context.Background()
	go func(ctx context.Context, port int) {
		hh.ginEngine.Run(":" + strconv.Itoa(port))
	}(ctx, hh.Port)
	return nil
}
func (hh *HttpApiServer) Install(env *core.XPluginEnv) (*core.XPluginMetaInfo, error) {
	return &core.XPluginMetaInfo{
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
// HttpApiServer Start
//
func (hh *HttpApiServer) Start(env *core.XPluginEnv) error {
	//
	// Render dashboard index
	//
	hh.ginEngine.GET(DASHBOARD_ROOT, hh.addRoute(Index))
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
	hh.ginEngine.POST(API_ROOT+"auth", hh.addRoute(Auth))
	//
	// Create InEnd
	//
	hh.ginEngine.POST(API_ROOT+"inends", hh.addRoute(InEnds))
	//
	// Create OutEnd
	//
	hh.ginEngine.POST(API_ROOT+"outends", hh.addRoute(OutEnds))
	//
	// Create rule
	//
	hh.ginEngine.POST(API_ROOT+"rules", hh.addRoute(Rules))
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
	log.Info("Http web dashboard started on:http://127.0.0.1:2580" + DASHBOARD_ROOT)
	return nil
}

func (hh *HttpApiServer) Uninstall(env *core.XPluginEnv) error {
	return nil
}
func (hh *HttpApiServer) Clean() {
}
