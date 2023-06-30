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

func (hs *HttpApiServer) Init(config *ini.Section) error {
	gin.SetMode(gin.ReleaseMode)
	hs.ginEngine = gin.New()

	var mainConfig _serverConfig
	if err := utils.InIMapToStruct(config, &mainConfig); err != nil {
		return err
	}
	hs.Host = mainConfig.Host
	hs.dbPath = mainConfig.DbPath
	hs.Port = mainConfig.Port
	hs.configHttpServer()
	//
	// Http server
	//
	go func(ctx context.Context, port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
		if err := hs.ginEngine.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("httpserver listen error: %s\n", err)
		}
	}(typex.GCTX, hs.Port)
	//
	// WebSocket server
	//
	hs.ginEngine.GET("/ws", glogger.WsLogger)
	return nil
}

/*
*
* 加载路由
*
 */
func (hs *HttpApiServer) LoadRoute() {
	//
	// Get all plugins
	//
	hs.ginEngine.GET(url("plugins"), hs.addRoute(Plugins))
	//
	// Get system information
	//
	hs.ginEngine.GET(url("system"), hs.addRoute(System))
	//
	// Ping -> Pong
	//
	hs.ginEngine.GET(url("ping"), hs.addRoute(Ping))
	//
	//
	//
	hs.ginEngine.GET(url("sourceCount"), hs.addRoute(SourceCount))
	//
	//
	//
	hs.ginEngine.GET(url("logs"), hs.addRoute(Logs))
	//
	//
	//
	hs.ginEngine.POST(url("logout"), hs.addRoute(LogOut))
	//
	// Get all inends
	//
	hs.ginEngine.GET(url("inends"), hs.addRoute(InEnds))
	hs.ginEngine.GET(url("inends/detail"), hs.addRoute(InEndDetail))
	//
	//
	//
	hs.ginEngine.GET(url("drivers"), hs.addRoute(Drivers))
	//
	// Get all outends
	//
	hs.ginEngine.GET(url("outends"), hs.addRoute(OutEnds))
	hs.ginEngine.GET(url("outends/detail"), hs.addRoute(OutEndDetail))
	//
	// Get all rules
	//
	hs.ginEngine.GET(url("rules"), hs.addRoute(Rules))
	hs.ginEngine.GET(url("rules/detail"), hs.addRoute(RuleDetail))
	//
	// Get statistics data
	//
	hs.ginEngine.GET(url("statistics"), hs.addRoute(Statistics))
	hs.ginEngine.GET(url("snapshot"), hs.addRoute(SnapshotDump))
	//
	// Auth
	//
	hs.ginEngine.GET(url("users"), hs.addRoute(Users))
	hs.ginEngine.GET(url("users/detail"), hs.addRoute(UserDetail))
	hs.ginEngine.POST(url("users"), hs.addRoute(CreateUser))
	//
	//
	//
	hs.ginEngine.POST(url("login"), hs.addRoute(Login))
	//
	//
	//
	hs.ginEngine.GET(url("info"), hs.addRoute(Info))
	//
	// Create InEnd
	//
	hs.ginEngine.POST(url("inends"), hs.addRoute(CreateInend))
	//
	// Update Inend
	//
	hs.ginEngine.PUT(url("inends"), hs.addRoute(UpdateInend))
	//
	// 配置表
	//
	hs.ginEngine.GET(url("inends/config"), hs.addRoute(GetInEndConfig))
	//
	// 数据模型表
	//
	hs.ginEngine.GET(url("inends/models"), hs.addRoute(GetInEndModels))
	//
	// Create OutEnd
	//
	hs.ginEngine.POST(url("outends"), hs.addRoute(CreateOutEnd))
	//
	// Update OutEnd
	//
	hs.ginEngine.PUT(url("outends"), hs.addRoute(UpdateOutEnd))
	//
	// Create rule
	//
	hs.ginEngine.POST(url("rules"), hs.addRoute(CreateRule))
	//
	// Update rule
	//
	hs.ginEngine.PUT(url("rules"), hs.addRoute(UpdateRule))
	//
	// Delete rule by UUID
	//
	hs.ginEngine.DELETE(url("rules"), hs.addRoute(DeleteRule))
	//
	// 测试规则
	//
	hs.ginEngine.POST(url("rules/testIn"), hs.addRoute(TestSourceCallback))
	hs.ginEngine.POST(url("rules/testOut"), hs.addRoute(TestOutEndCallback))
	hs.ginEngine.POST(url("rules/testDevice"), hs.addRoute(TestDeviceCallback))
	//
	// Delete inend by UUID
	//
	hs.ginEngine.DELETE(url("inends"), hs.addRoute(DeleteInEnd))
	//
	// Delete outEnd by UUID
	//
	hs.ginEngine.DELETE(url("outends"), hs.addRoute(DeleteOutEnd))

	//
	// 验证 lua 语法
	//
	hs.ginEngine.POST(url("validateRule"), hs.addRoute(ValidateLuaSyntax))
	//
	// 获取配置表
	//
	hs.ginEngine.GET(url("rType"), hs.addRoute(RType))
	hs.ginEngine.GET(url("tType"), hs.addRoute(TType))
	hs.ginEngine.GET(url("dType"), hs.addRoute(DType))
	//
	// 串口列表
	//
	hs.ginEngine.GET(url("uarts"), hs.addRoute(GetUarts))
	//
	// 获取服务启动时间
	//
	hs.ginEngine.GET(url("startedAt"), hs.addRoute(StartedAt))
	//
	// 设备管理
	//
	hs.ginEngine.GET(url("devices"), hs.addRoute(Devices))
	hs.ginEngine.GET(url("devices/detail"), hs.addRoute(DeviceDetail))
	hs.ginEngine.POST(url("devices"), hs.addRoute(CreateDevice))
	hs.ginEngine.PUT(url("devices"), hs.addRoute(UpdateDevice))
	hs.ginEngine.DELETE(url("devices"), hs.addRoute(DeleteDevice))
	hs.ginEngine.POST(url("devices/modbus/sheetImport"), hs.addRoute(ModbusSheetImport))

	// 外挂管理
	hs.ginEngine.GET(url("goods"), hs.addRoute(Goods))
	hs.ginEngine.POST(url("goods"), hs.addRoute(CreateGoods))
	hs.ginEngine.PUT(url("goods"), hs.addRoute(UpdateGoods))
	hs.ginEngine.DELETE(url("goods"), hs.addRoute(DeleteGoods))
	// 加载资源类型
	source.LoadSt()
	target.LoadTt()
	device.LoadDt()
	// ----------------------------------------------------------------------------------------------
	// APP
	// ----------------------------------------------------------------------------------------------
	hs.ginEngine.GET(url("app"), hs.addRoute(Apps))
	hs.ginEngine.POST(url("app"), hs.addRoute(CreateApp))
	hs.ginEngine.PUT(url("app"), hs.addRoute(UpdateApp))
	hs.ginEngine.DELETE(url("app"), hs.addRoute(RemoveApp))
	hs.ginEngine.PUT(url("app/start"), hs.addRoute(StartApp))
	hs.ginEngine.PUT(url("app/stop"), hs.addRoute(StopApp))
	hs.ginEngine.GET(url("app/detail"), hs.addRoute(AppDetail))
	// ----------------------------------------------------------------------------------------------
	// AI BASE
	// ----------------------------------------------------------------------------------------------
	hs.ginEngine.GET(url("aibase"), hs.addRoute(AiBase))
	hs.ginEngine.DELETE(url("aibase"), hs.addRoute(DeleteAiBase))
	// ----------------------------------------------------------------------------------------------
	// Plugin
	// ----------------------------------------------------------------------------------------------
	hs.ginEngine.POST(url("plugin/service"), hs.addRoute(PluginService))
	hs.ginEngine.GET(url("plugin/detail"), hs.addRoute(PluginDetail))

}

// HttpApiServer Start
func (hs *HttpApiServer) Start(r typex.RuleX) error {
	hs.ruleEngine = r
	hs.LoadRoute()
	glogger.GLogger.Infof("Http server started on http://0.0.0.0:%v", hs.Port)
	return nil
}

func (hs *HttpApiServer) Stop() error {
	return nil
}

func (hs *HttpApiServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hs.uuid,
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
func (*HttpApiServer) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{Out: "HttpApiServer"}
}

// Add api route
func (hs *HttpApiServer) addRoute(f func(*gin.Context, *HttpApiServer)) func(*gin.Context) {

	return func(c *gin.Context) {
		f(c, hs)
	}
}
func (hs *HttpApiServer) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (hs *HttpApiServer) configHttpServer() {
	hs.ginEngine.Use(hs.Authorize())
	hs.ginEngine.Use(common.Cros())
	hs.ginEngine.Use(static.Serve("/", WWWRoot("")))
	hs.ginEngine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		glogger.GLogger.Error(err)
		c.JSON(common.HTTP_OK, common.Error500(err1crash))
	}))
	hs.ginEngine.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})
	if hs.dbPath == "" {
		hs.InitDb(_DEFAULT_DB_PATH)
	} else {
		hs.InitDb(hs.dbPath)
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
