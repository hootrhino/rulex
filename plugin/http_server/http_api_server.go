package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hootrhino/rulex/core"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"

	"github.com/gin-contrib/static"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	uuid       string
}

/*
*
* 初始化数据库
*
 */
func (s *HttpApiServer) InitDb(dbPath string) {
	var err error
	if core.GlobalConfig.AppDebugMode {
		s.sqliteDb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		s.sqliteDb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}

	if err != nil {
		glogger.GLogger.Error(err)
		// Sqlite 创建失败应该是致命错误了, 多半是环境出问题，直接给panic了, 不尝试救活
		panic(err)
	}
	// 注册数据库配置表
	// 这么写看起来是很难受, 但是这玩意就是go的哲学啊(大道至简？？？)
	if err := s.DB().AutoMigrate(
		&model.MInEnd{},
		&model.MOutEnd{},
		&model.MRule{},
		&model.MUser{},
		&model.MDevice{},
		&model.MGoods{},
		&model.MApp{},
		&model.MAiBase{},
		&model.MModbusPointPosition{},
	); err != nil {
		glogger.GLogger.Fatal(err)
		os.Exit(1)
	}
}

func (s *HttpApiServer) DB() *gorm.DB {
	return s.sqliteDb
}
func NewHttpApiServer() *HttpApiServer {
	return &HttpApiServer{
		uuid: "HTTP-API-SERVER",
	}
}

// HTTP服务器崩了, 重启恢复
var err1crash = errors.New("http server crash, try to recovery")

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
	hh.configHttpServer()
	//
	// Http server
	//
	go func(ctx context.Context, port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
		if err := hh.ginEngine.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
	}(typex.GCTX, hh.Port)
	//
	// WebSocket server
	//
	hh.ginEngine.GET("/ws", glogger.WsLogger)
	return nil
}
func (hh *HttpApiServer) LoadRoute() {
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
	hh.ginEngine.GET(url("inends/detail"), hh.addRoute(InEndDetail))
	//
	//
	//
	hh.ginEngine.GET(url("drivers"), hh.addRoute(Drivers))
	//
	// Get all outends
	//
	hh.ginEngine.GET(url("outends"), hh.addRoute(OutEnds))
	hh.ginEngine.GET(url("outends/detail"), hh.addRoute(OutEndDetail))
	//
	// Get all rules
	//
	hh.ginEngine.GET(url("rules"), hh.addRoute(Rules))
	hh.ginEngine.GET(url("rules/detail"), hh.addRoute(RuleDetail))
	//
	// Get statistics data
	//
	hh.ginEngine.GET(url("statistics"), hh.addRoute(Statistics))
	hh.ginEngine.GET(url("snapshot"), hh.addRoute(SnapshotDump))
	//
	// Auth
	//
	hh.ginEngine.GET(url("users"), hh.addRoute(Users))
	hh.ginEngine.GET(url("users/detail"), hh.addRoute(UserDetail))
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
	//
	// Update Inend
	//
	hh.ginEngine.PUT(url("inends"), hh.addRoute(UpdateInend))
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
	//
	// Update OutEnd
	//
	hh.ginEngine.PUT(url("outends"), hh.addRoute(UpdateOutEnd))
	//
	// Create rule
	//
	hh.ginEngine.POST(url("rules"), hh.addRoute(CreateRule))
	//
	// Update rule
	//
	hh.ginEngine.PUT(url("rules"), hh.addRoute(UpdateRule))
	//
	// Delete rule by UUID
	//
	hh.ginEngine.DELETE(url("rules"), hh.addRoute(DeleteRule))
	hh.ginEngine.POST(url("rules/testIn"), hh.addRoute(TestSourceCallback))
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
	hh.ginEngine.GET(url("devices/detail"), hh.addRoute(DeviceDetail))
	hh.ginEngine.POST(url("devices"), hh.addRoute(CreateDevice))
	hh.ginEngine.PUT(url("devices"), hh.addRoute(UpdateDevice))
	hh.ginEngine.DELETE(url("devices"), hh.addRoute(DeleteDevice))
	hh.ginEngine.POST(url("devices/modbus/sheetImport"), hh.addRoute(ModbusSheetImport))

	// 外挂管理
	hh.ginEngine.GET(url("goods"), hh.addRoute(Goods))
	hh.ginEngine.POST(url("goods"), hh.addRoute(CreateGoods))
	hh.ginEngine.PUT(url("goods"), hh.addRoute(UpdateGoods))
	hh.ginEngine.DELETE(url("goods"), hh.addRoute(DeleteGoods))
	// 加载资源类型
	source.LoadSt()
	target.LoadTt()
	device.LoadDt()
	// ----------------------------------------------------------------------------------------------
	// APP
	// ----------------------------------------------------------------------------------------------
	hh.ginEngine.GET(url("app"), hh.addRoute(Apps))
	hh.ginEngine.POST(url("app"), hh.addRoute(CreateApp))
	hh.ginEngine.PUT(url("app"), hh.addRoute(UpdateApp))
	hh.ginEngine.DELETE(url("app"), hh.addRoute(RemoveApp))
	hh.ginEngine.PUT(url("app/start"), hh.addRoute(StartApp))
	hh.ginEngine.PUT(url("app/stop"), hh.addRoute(StopApp))
	hh.ginEngine.GET(url("app/detail"), hh.addRoute(AppDetail))
	// ----------------------------------------------------------------------------------------------
	// AI BASE
	// ----------------------------------------------------------------------------------------------
	hh.ginEngine.GET(url("aibase"), hh.addRoute(AiBase))
	hh.ginEngine.DELETE(url("aibase"), hh.addRoute(DeleteAiBase))
	// ----------------------------------------------------------------------------------------------
	// Plugin
	// ----------------------------------------------------------------------------------------------
	hh.ginEngine.POST(url("plugin/service"), hh.addRoute(PluginService))
	hh.ginEngine.GET(url("plugin/detail"), hh.addRoute(PluginDetail))

}

// HttpApiServer Start
func (hh *HttpApiServer) Start(r typex.RuleX) error {
	hh.ruleEngine = r
	hh.LoadRoute()
	glogger.GLogger.Infof("Http server started on http://0.0.0.0:%v", hh.Port)
	return nil
}

func (hh *HttpApiServer) Stop() error {
	return nil
}

func (hh *HttpApiServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "RULEX HTTP RESTFul Api Server",
		Version:  "v1.0.0",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
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
func (cs *HttpApiServer) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{Out: "HttpApiServer"}
}

// Add api route
func (h *HttpApiServer) addRoute(f func(*gin.Context, *HttpApiServer)) func(*gin.Context) {

	return func(c *gin.Context) {
		f(c, h)
	}
}
func (hh *HttpApiServer) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (hh *HttpApiServer) configHttpServer() {
	hh.ginEngine.Use(hh.Authorize())
	hh.ginEngine.Use(common.Cros())
	hh.ginEngine.Use(static.Serve("/", WWWRoot("")))
	hh.ginEngine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		glogger.GLogger.Error(err)
		c.JSON(common.HTTP_OK, common.Error500(err1crash))
	}))
	hh.ginEngine.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})
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
