package httpserver

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-contrib/static"
	"github.com/i4de/rulex/device"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/source"
	"github.com/i4de/rulex/target"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

const _API_V1_ROOT string = "/api/v1/"
const _DEFAULT_DB_PATH string = "./rulex.db"

// 启动时间
var StartedTime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	DbPath string `ini:"dbpath"`
	Port   int    `ini:"port"`
}
type HttpApiServer struct {
	Port       int
	Host       string
	sqliteDb   *gorm.DB
	dbPath     string
	ginEngine  *gin.Engine
	ruleEngine typex.RuleX
}

func NewHttpApiServer() *HttpApiServer {
	return &HttpApiServer{}
}

func (hh *HttpApiServer) Init(config *ini.Section) error {
	gin.SetMode(gin.ReleaseMode)
	hh.ginEngine = gin.New()
	var mainConfig _serverConfig
	if err := utils.InIMapToStruct(config, &mainConfig); err != nil {
		return err
	}
	hh.Host = mainConfig.Host
	hh.dbPath = mainConfig.DbPath
	hh.Port = mainConfig.Port
	configHttpServer(hh)
	//
	// Http server
	//
	go func(ctx context.Context, port int) {
		if err := hh.ginEngine.Run(":" + strconv.Itoa(port)); err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
	}(typex.GCTX, hh.Port)
	//
	// WebSocket server
	//
	hh.ginEngine.GET("/ws", glogger.WsLogger)
	return nil
}

// HttpApiServer Start
func (hh *HttpApiServer) Start(r typex.RuleX) error {
	hh.ruleEngine = r
	//
	// Get all plugins
	//
	hh.ginEngine.GET(url("plugins"), hh.addRoute(Plugins))
	//
	// Get system information
	//
	hh.ginEngine.GET(url("system"), hh.addRoute(System))
	//
	// Ping -> Pong
	//
	hh.ginEngine.GET(url("ping"), hh.addRoute(Ping))
	//
	//
	//
	hh.ginEngine.GET(url("sourceCount"), hh.addRoute(SourceCount))
	//
	//
	//
	hh.ginEngine.GET(url("logs"), hh.addRoute(Logs))
	//
	//
	//
	hh.ginEngine.POST(url("logout"), hh.addRoute(LogOut))
	//
	// Get all inends
	//
	hh.ginEngine.GET(url("inends"), hh.addRoute(InEnds))
	//
	//
	//
	hh.ginEngine.GET(url("drivers"), hh.addRoute(Drivers))
	//
	// Get all outends
	//
	hh.ginEngine.GET(url("outends"), hh.addRoute(OutEnds))
	//
	// Get all rules
	//
	hh.ginEngine.GET(url("rules"), hh.addRoute(Rules))
	//
	// Get statistics data
	//
	hh.ginEngine.GET(url("statistics"), hh.addRoute(Statistics))
	hh.ginEngine.GET(url("snapshot"), hh.addRoute(SnapshotDump))
	//
	// Auth
	//
	hh.ginEngine.GET(url("users"), hh.addRoute(Users))
	hh.ginEngine.POST(url("users"), hh.addRoute(CreateUser))
	//
	//
	//
	hh.ginEngine.POST(url("login"), hh.addRoute(Login))
	//
	//
	//
	hh.ginEngine.GET(url("info"), hh.addRoute(Info))
	//
	// Create InEnd
	//
	hh.ginEngine.POST(url("inends"), hh.addRoute(CreateInend))
	hh.ginEngine.PUT(url("inends"), hh.addRoute(CreateInend))
	//
	// 配置表
	//
	hh.ginEngine.GET(url("inends/config"), hh.addRoute(GetInEndConfig))
	//
	// 数据模型表
	//
	hh.ginEngine.GET(url("inends/models"), hh.addRoute(GetInEndModels))
	//
	// Create OutEnd
	//
	hh.ginEngine.POST(url("outends"), hh.addRoute(CreateOutEnd))
	hh.ginEngine.PUT(url("outends"), hh.addRoute(CreateOutEnd))
	//
	// Create rule
	//
	hh.ginEngine.POST(url("rules"), hh.addRoute(CreateRule))
	//
	// Update rule
	//
	hh.ginEngine.PUT(url("rules"), hh.addRoute(CreateRule))
	//
	// Delete rule by UUID
	//
	hh.ginEngine.DELETE(url("rules"), hh.addRoute(DeleteRule))
	//
	// Delete inend by UUID
	//
	hh.ginEngine.DELETE(url("inends"), hh.addRoute(DeleteInEnd))
	//
	// Delete outEnd by UUID
	//
	hh.ginEngine.DELETE(url("outends"), hh.addRoute(DeleteOutEnd))

	//
	// 验证 lua 语法
	//
	hh.ginEngine.POST(url("validateRule"), hh.addRoute(ValidateLuaSyntax))
	//
	// 获取配置表
	//
	hh.ginEngine.GET(url("rType"), hh.addRoute(RType))
	hh.ginEngine.GET(url("tType"), hh.addRoute(TType))
	hh.ginEngine.GET(url("dType"), hh.addRoute(DType))
	//
	// 串口列表
	//
	hh.ginEngine.GET(url("uarts"), hh.addRoute(GetUarts))
	//
	// 获取服务启动时间
	//
	hh.ginEngine.GET(url("startedAt"), hh.addRoute(StartedAt))
	//
	// 设备管理
	//
	hh.ginEngine.GET(url("devices"), hh.addRoute(Devices))
	hh.ginEngine.POST(url("devices"), hh.addRoute(CreateDevice))
	hh.ginEngine.PUT(url("devices"), hh.addRoute(UpdateDevice))
	hh.ginEngine.DELETE(url("devices"), hh.addRoute(DeleteDevice))

	// 外挂管理
	hh.ginEngine.GET(url("goods"), hh.addRoute(Goods))
	hh.ginEngine.POST(url("goods"), hh.addRoute(CreateGoods))
	hh.ginEngine.PUT(url("goods"), hh.addRoute(UpdateGoods))
	hh.ginEngine.DELETE(url("goods"), hh.addRoute(DeleteGoods))
	// 加载资源类型
	source.LoadSt()
	target.LoadTt()
	device.LoadDt()
	glogger.GLogger.Infof("Http server started on http://0.0.0.0:%v", hh.Port)
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
		Name:     "Rulex Base Api Server",
		Version:  typex.DefaultVersion.Version,
		Homepage: "https://rulex.pages.dev",
		HelpLink: "https://rulex.pages.dev",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

/*
*
* 服务调用接口
*
 */
func (cs *HttpApiServer) Service(arg typex.ServiceArg) error {
	return nil
}
//--------------------------------------------------------------------------------
//go:embed  www/*
var files embed.FS

type WWWFS struct {
	http.FileSystem
}

func (f WWWFS) Exists(prefix string, filepath string) bool {
	_, err := f.Open(path.Join(prefix, filepath))
	return err == nil
}

func wwwRoot(dir string) static.ServeFileSystem {
	if sub, err := fs.Sub(files, path.Join("www", dir)); err == nil {
		return WWWFS{http.FS(sub)}
	}
	return nil
}

func configHttpServer(hh *HttpApiServer) {
	hh.ginEngine.Use(hh.Authorize())
	hh.ginEngine.Use(Cros())
	hh.ginEngine.Use(static.Serve("/", wwwRoot("")))

	if hh.dbPath == "" {
		hh.InitDb(_DEFAULT_DB_PATH)
	} else {
		hh.InitDb(hh.dbPath)
	}
}

/*
*
* 拼接URL
*
 */
func url(path string) string {
	return _API_V1_ROOT + path
}
