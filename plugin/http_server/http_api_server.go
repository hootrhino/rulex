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
	ruleEngine core.RuleX
}

func NewHttpApiServer(port int, root string, e core.RuleX) *HttpApiServer {
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
	hh.ginEngine.LoadHTMLFiles(hh.Root+"/login.html", hh.Root+"/view/rulex/index.html")
	hh.ginEngine.Static("/dashboard/v1/component", hh.Root+"/component")
	hh.ginEngine.Static("/dashboard/v1/admin", hh.Root+"/admin")
	hh.ginEngine.Static("/dashboard/v1/view", hh.Root+"/view")
	hh.ginEngine.Static("/dashboard/v1/config", hh.Root+"/config")
	hh.ginEngine.Static("/component", hh.Root+"/component")
	hh.ginEngine.Static("/admin", hh.Root+"/admin")
	hh.ginEngine.Static("/view", hh.Root+"/view")
	hh.ginEngine.Static("/config", hh.Root+"/config")

	// 后台管理员页面
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
	hh.ginEngine.GET("/", hh.addRoute(Login))
	hh.ginEngine.GET(DASHBOARD_ROOT, hh.addRoute(Login))
	hh.ginEngine.GET(DASHBOARD_ROOT+"login", hh.addRoute(Login))
	hh.ginEngine.GET(DASHBOARD_ROOT+"index", hh.addRoute(Index))
	//
	// List CloudServices
	//
	hh.ginEngine.GET(API_ROOT+"cloudServices", hh.addRoute(CloudServices))

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
	log.Info("Http server started on http://127.0.0.1:2580" + DASHBOARD_ROOT)
	return nil
}

func (hh *HttpApiServer) Uninstall(env *core.XPluginEnv) error {
	return nil
}
func (hh *HttpApiServer) Clean() {
}
