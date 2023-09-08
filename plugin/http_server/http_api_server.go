package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	common "github.com/hootrhino/rulex/plugin/http_server/common"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"

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
	uuid       string
	ginEngine  *gin.Engine
	ruleEngine typex.RuleX
	mainConfig _serverConfig
}

/*
*
* 初始化数据库
*
 */
func (s *HttpApiServer) registerModel() {
	sqlitedao.Sqlite.DB().AutoMigrate(
		&model.MInEnd{},
		&model.MOutEnd{},
		&model.MRule{},
		&model.MUser{},
		&model.MDevice{},
		&model.MGoods{},
		&model.MApp{},
		&model.MAiBase{},
		&model.MModbusPointPosition{},
		&model.MVisual{},
		&model.MGenericGroup{},
		&model.MGenericGroupRelation{},
		&model.MProtocolApp{},
		&model.MNetworkConfig{},
		&model.MWifiConfig{},
	)
}

func NewHttpApiServer() *HttpApiServer {
	return &HttpApiServer{
		uuid:       "HTTP-API-SERVER",
		mainConfig: _serverConfig{},
	}
}

// HTTP服务器崩了, 重启恢复
var err1crash = errors.New("http server crash, try to recovery")

func (hs *HttpApiServer) Init(config *ini.Section) error {
	gin.SetMode(gin.ReleaseMode)
	hs.ginEngine = gin.New()
	if err := utils.InIMapToStruct(config, &hs.mainConfig); err != nil {
		return err
	}
	if hs.mainConfig.DbPath == "" {
		sqlitedao.Load(_DEFAULT_DB_PATH)

	} else {
		sqlitedao.Load(hs.mainConfig.DbPath)
	}
	hs.registerModel()
	hs.configHttpServer()
	hs.InitializeNetWorkConfigData()
	//
	// WebSocket server
	//
	hs.ginEngine.GET("/ws", glogger.WsLogger)
	return nil
}

/*
*
* 初始化网络配置
*
 */
func (hs *HttpApiServer) InitializeNetWorkConfigData() {
	// 初始化有线网口配置
	if !service.CheckIfAlreadyInitNetWorkConfig() {
		service.InitNetWorkConfig()
	}
	// 初始化WIFI配置
	if !service.CheckIfAlreadyInitWlanConfig() {
		service.InitWlanConfig()
	}
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
	// 网络适配器列表
	//
	osApi := hs.ginEngine.Group(url("/os"))
	{
		osApi.GET(("/netInterfaces"), hs.addRoute(GetNetInterfaces))
		osApi.GET(("/osRelease"), hs.addRoute(CatOsRelease))
		osApi.GET(("/uarts"), hs.addRoute(GetUartList))
		osApi.GET(("/system"), hs.addRoute(System))
		osApi.GET(("/startedAt"), hs.addRoute(StartedAt))

	}
	//
	// 设备管理
	//
	deviceApi := hs.ginEngine.Group(url("/devices"))
	{
		deviceApi.GET(("/"), hs.addRoute(Devices))
		deviceApi.POST(("/"), hs.addRoute(CreateDevice))
		deviceApi.PUT(("/"), hs.addRoute(UpdateDevice))
		deviceApi.DELETE(("/"), hs.addRoute(DeleteDevice))
		deviceApi.GET(("/detail"), hs.addRoute(DeviceDetail))
		deviceApi.POST(("/modbus/sheetImport"), hs.addRoute(ModbusSheetImport))
		deviceApi.PUT(("/modbus/point"), hs.addRoute(UpdateModbusPoint))
		deviceApi.GET(("/modbus"), hs.addRoute(ModbusPoints))
	}
	goodsApi := hs.ginEngine.Group(url("/goods"))
	{
		// 外挂管理
		goodsApi.GET(("/"), hs.addRoute(Goods))
		goodsApi.POST(("/"), hs.addRoute(CreateGoods))
		goodsApi.PUT(("/"), hs.addRoute(UpdateGoods))
		goodsApi.DELETE(("/"), hs.addRoute(DeleteGoods))
	}

	// 加载资源类型
	source.LoadSt()
	target.LoadTt()
	device.LoadDt()
	// ----------------------------------------------------------------------------------------------
	// APP
	// ----------------------------------------------------------------------------------------------
	appApi := hs.ginEngine.Group(url("/app"))
	{
		appApi.GET(("/"), hs.addRoute(Apps))
		appApi.POST(("/"), hs.addRoute(CreateApp))
		appApi.PUT(("/"), hs.addRoute(UpdateApp))
		appApi.DELETE(("/"), hs.addRoute(RemoveApp))
		appApi.PUT(("/start"), hs.addRoute(StartApp))
		appApi.PUT(("/stop"), hs.addRoute(StopApp))
		appApi.GET(("/detail"), hs.addRoute(AppDetail))
	}
	// ----------------------------------------------------------------------------------------------
	// AI BASE
	// ----------------------------------------------------------------------------------------------
	aiApi := hs.ginEngine.Group(url("/aibase"))
	{
		aiApi.GET(("/"), hs.addRoute(AiBase))
		aiApi.DELETE(("/"), hs.addRoute(DeleteAiBase))
	}
	// ----------------------------------------------------------------------------------------------
	// Plugin
	// ----------------------------------------------------------------------------------------------
	pluginApi := hs.ginEngine.Group(url("/plugin"))
	{
		pluginApi.POST(("/service"), hs.addRoute(PluginService))
		pluginApi.GET(("/detail"), hs.addRoute(PluginDetail))
	}

	//
	// 分组管理
	//
	groupApi := hs.ginEngine.Group(url("/group"))
	{
		groupApi.POST("/create", hs.addRoute(CreateGroup))
		groupApi.DELETE("/delete", hs.addRoute(DeleteGroup))
		groupApi.PUT("/update", hs.addRoute(UpdateGroup))
		groupApi.GET("/list", hs.addRoute(ListGroup))
		groupApi.POST("/bind", hs.addRoute(BindResource))
		groupApi.PUT("/unbind", hs.addRoute(UnBindResource))
		groupApi.GET("/devices", hs.addRoute(FindDeviceByGroup))
		groupApi.GET("/visuals", hs.addRoute(FindVisualByGroup))
	}

	//
	// 协议应用管理
	//
	protoAppApi := hs.ginEngine.Group(url("/protoapp"))
	{
		protoAppApi.POST("/create", hs.addRoute(CreateProtocolApp))
		protoAppApi.DELETE("/delete", hs.addRoute(DeleteProtocolApp))
		protoAppApi.PUT("/update", hs.addRoute(UpdateProtocolApp))
		protoAppApi.GET("/list", hs.addRoute(ListProtocolApp))
	}
	//
	// 大屏应用管理
	//
	screenApi := hs.ginEngine.Group(url("/visual"))
	{
		screenApi.POST("/create", hs.addRoute(CreateVisual))
		screenApi.DELETE("/delete", hs.addRoute(DeleteVisual))
		screenApi.PUT("/update", hs.addRoute(UpdateVisual))
		screenApi.GET("/list", hs.addRoute(ListVisual))
	}
	//
	// 系统设置
	//
	hs.LoadSystemSettingsAPI()
}

// HttpApiServer Start
func (hs *HttpApiServer) Start(r typex.RuleX) error {
	hs.ruleEngine = r
	hs.LoadRoute()
	glogger.GLogger.Infof("Http server started on :%v", hs.mainConfig.Port)
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
	}(typex.GCTX, hs.mainConfig.Port)
}

/*
*
* 拼接URL
*
 */
func url(path string) string {
	return _API_V1_ROOT + path
}
