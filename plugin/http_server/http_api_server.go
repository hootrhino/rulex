package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/hootrhino/rulex/component/cron_task"

	"github.com/hootrhino/rulex/component/appstack"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/plugin/http_server/apis"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/server"
	"github.com/hootrhino/rulex/plugin/http_server/service"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"gopkg.in/ini.v1"

	_ "github.com/mattn/go-sqlite3"
)

type _serverConfig struct {
	DbPath string `ini:"dbpath"`
}
type ApiServerPlugin struct {
	uuid       string
	ruleEngine typex.RuleX
	mainConfig _serverConfig
}

func NewHttpApiServer(ruleEngine typex.RuleX) *ApiServerPlugin {
	return &ApiServerPlugin{
		uuid:       "HTTP-API-SERVER",
		mainConfig: _serverConfig{},
		ruleEngine: ruleEngine,
	}
}

/*
*
* 初始化RULEX
*
 */
func initRulex(engine typex.RuleX) {
	/*
	*
	* 加载schema到内存中
	*
	 */
	for _, mDataSchema := range service.AllDataSchema() {
		dataDefine := []typex.DataDefine{}
		err := json.Unmarshal([]byte(mDataSchema.Schema), &dataDefine)
		if err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		// 初始化装入ne
		core.SchemaSet(mDataSchema.UUID, typex.DataSchema{
			UUID:   mDataSchema.UUID,
			Name:   mDataSchema.Name,
			Type:   mDataSchema.Type,
			Schema: dataDefine,
		})
	}
	//
	// Load inend from sqlite
	//
	for _, minEnd := range service.AllMInEnd() {
		if err := server.LoadNewestInEnd(minEnd.UUID, engine); err != nil {
			glogger.GLogger.Error("InEnd load failed:", err)
		}
	}

	//
	// Load out from sqlite
	//
	for _, mOutEnd := range service.AllMOutEnd() {
		if err := server.LoadNewestOutEnd(mOutEnd.UUID, engine); err != nil {
			glogger.GLogger.Error("OutEnd load failed:", err)
		}
	}
	// 加载设备
	for _, mDevice := range service.AllDevices() {
		glogger.GLogger.Debug("LoadNewestDevice mDevice.BindRules: ", mDevice.BindRules.String())
		if err := server.LoadNewestDevice(mDevice.UUID, engine); err != nil {
			glogger.GLogger.Error("Device load failed:", err)
		}

	}
	// 加载外挂
	for _, mGoods := range service.AllGoods() {
		newGoods := trailer.Goods{
			UUID:        mGoods.UUID,
			LocalPath:   mGoods.LocalPath,
			NetAddr:     mGoods.NetAddr,
			Description: mGoods.Description,
			Args:        mGoods.Args,
		}
		if err := trailer.Fork(newGoods); err != nil {
			glogger.GLogger.Error("Goods load failed:", err)
		}
	}
	//
	// APP stack
	//
	for _, mApp := range service.AllApp() {
		app := appstack.NewApplication(
			mApp.UUID,
			mApp.Name,
			mApp.Version,
		)
		if err := appstack.LoadApp(app, mApp.LuaSource); err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		if *mApp.AutoStart {
			glogger.GLogger.Debug("App autoStart allowed:", app.UUID, app.Version, app.Name)
			if err1 := appstack.StartApp(app.UUID); err1 != nil {
				glogger.GLogger.Error("App autoStart failed:", err1)
			}
		}
	}
}
func (hs *ApiServerPlugin) Init(config *ini.Section) error {
	if err := utils.InIMapToStruct(config, &hs.mainConfig); err != nil {
		return err
	}
	server.StartRulexApiServer(hs.ruleEngine)
	interdb.RegisterModel(
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
		&model.MDataSchema{},
		&model.MSiteConfig{},
		&model.MIpRoute{},
		&model.MCronTask{},
		&model.MCronResult{},
	)
	server.DefaultApiServer.InitializeData()
	initRulex(hs.ruleEngine)
	return nil
}

/*
*
* 初始化网络配置
*
 */
func (hs *ApiServerPlugin) InitializeData() {
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
}

/*
*
* 加载路由
*
 */
func (hs *ApiServerPlugin) LoadRoute() {
	//
	// Get all plugins
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("plugins"), server.AddRoute(apis.Plugins))
	//
	// Get system information
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("system"), server.AddRoute(apis.System))
	//
	// Ping -> Pong
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("ping"), server.AddRoute(apis.Ping))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("sourceCount"), server.AddRoute(apis.SourceCount))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("logs"), server.AddRoute(apis.Logs))
	//
	//
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("logout"), server.AddRoute(apis.LogOut))
	//
	// Get all inends
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends"), server.AddRoute(apis.InEnds))
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends/detail"), server.AddRoute(apis.InEndDetail))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("drivers"), server.AddRoute(apis.Drivers))
	//
	// Get all outends
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("outends"), server.AddRoute(apis.OutEnds))
	server.DefaultApiServer.Route().GET(server.ContextUrl("outends/detail"), server.AddRoute(apis.OutEndDetail))
	//
	// Get all rules
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("rules"), server.AddRoute(apis.Rules))
	server.DefaultApiServer.Route().GET(server.ContextUrl("rules/detail"), server.AddRoute(apis.RuleDetail))
	//
	// Get statistics data
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("statistics"), server.AddRoute(apis.Statistics))
	server.DefaultApiServer.Route().GET(server.ContextUrl("snapshot"), server.AddRoute(apis.SnapshotDump))
	//
	// Auth
	//
	userApi := server.RouteGroup(server.ContextUrl("/users"))
	{
		userApi.GET(("/"), server.AddRoute(apis.Users))
		userApi.GET(("/detail"), server.AddRoute(apis.UserDetail))
		userApi.POST(("/"), server.AddRoute(apis.CreateUser))

	}

	//
	//
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("login"), server.AddRoute(apis.Login))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("info"), server.AddRoute(apis.Info))
	//
	// Create InEnd
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("inends"), server.AddRoute(apis.CreateInend))
	//
	// Update Inend
	//
	server.DefaultApiServer.Route().PUT(server.ContextUrl("inends"), server.AddRoute(apis.UpdateInend))
	//
	// 配置表
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends/config"), server.AddRoute(apis.GetInEndConfig))
	//
	// 数据模型表
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends/models"), server.AddRoute(apis.GetInEndModels))
	//
	// Create OutEnd
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("outends"), server.AddRoute(apis.CreateOutEnd))
	//
	// Update OutEnd
	//
	server.DefaultApiServer.Route().PUT(server.ContextUrl("outends"), server.AddRoute(apis.UpdateOutEnd))
	rulesApi := server.RouteGroup(server.ContextUrl("/rules"))
	{
		rulesApi.POST(("/"), server.AddRoute(apis.CreateRule))
		rulesApi.PUT(("/"), server.AddRoute(apis.UpdateRule))
		rulesApi.DELETE(("/"), server.AddRoute(apis.DeleteRule))
		rulesApi.POST(("/testIn"), server.AddRoute(apis.TestSourceCallback))
		rulesApi.POST(("/testOut"), server.AddRoute(apis.TestOutEndCallback))
		rulesApi.POST(("/testDevice"), server.AddRoute(apis.TestDeviceCallback))
	}

	//
	// Delete inend by UUID
	//
	server.DefaultApiServer.Route().DELETE(server.ContextUrl("inends"), server.AddRoute(apis.DeleteInEnd))
	//
	// Delete outEnd by UUID
	//
	server.DefaultApiServer.Route().DELETE(server.ContextUrl("outends"), server.AddRoute(apis.DeleteOutEnd))

	//
	// 验证 lua 语法
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("validateRule"), server.AddRoute(apis.ValidateLuaSyntax))
	//
	// 获取配置表
	//
	resourceTypeApi := server.RouteGroup(server.ContextUrl("/"))
	{
		resourceTypeApi.GET(("rType"), server.AddRoute(apis.RType))
		resourceTypeApi.GET(("tType"), server.AddRoute(apis.TType))
		resourceTypeApi.GET(("dType"), server.AddRoute(apis.DType))
	}
	//
	// 网络适配器列表
	//
	osApi := server.RouteGroup(server.ContextUrl("/os"))
	{
		osApi.GET(("/netInterfaces"), server.AddRoute(apis.GetNetInterfaces))
		osApi.GET(("/osRelease"), server.AddRoute(apis.CatOsRelease))
		osApi.GET(("/uarts"), server.AddRoute(apis.GetUartList))
		osApi.GET(("/system"), server.AddRoute(apis.System))
		osApi.GET(("/startedAt"), server.AddRoute(apis.StartedAt))

	}
	//
	// 设备管理
	//
	deviceApi := server.RouteGroup(server.ContextUrl("/devices"))
	{
		deviceApi.GET(("/"), server.AddRoute(apis.Devices))
		deviceApi.POST(("/"), server.AddRoute(apis.CreateDevice))
		deviceApi.PUT(("/"), server.AddRoute(apis.UpdateDevice))
		deviceApi.DELETE(("/"), server.AddRoute(apis.DeleteDevice))
		deviceApi.GET(("/detail"), server.AddRoute(apis.DeviceDetail))
		deviceApi.POST(("/modbus/sheetImport"), server.AddRoute(apis.ModbusSheetImport))
		deviceApi.PUT(("/modbus/point"), server.AddRoute(apis.UpdateModbusPoint))
		deviceApi.GET(("/modbus"), server.AddRoute(apis.ModbusPoints))
		deviceApi.GET("/group", server.AddRoute(apis.ListDeviceGroup))

	}

	// ----------------------------------------------------------------------------------------------
	// APP
	// ----------------------------------------------------------------------------------------------
	appApi := server.RouteGroup(server.ContextUrl("/app"))
	{
		appApi.GET(("/"), server.AddRoute(apis.Apps))
		appApi.POST(("/"), server.AddRoute(apis.CreateApp))
		appApi.PUT(("/"), server.AddRoute(apis.UpdateApp))
		appApi.DELETE(("/"), server.AddRoute(apis.RemoveApp))
		appApi.PUT(("/start"), server.AddRoute(apis.StartApp))
		appApi.PUT(("/stop"), server.AddRoute(apis.StopApp))
		appApi.GET(("/detail"), server.AddRoute(apis.AppDetail))
	}
	// ----------------------------------------------------------------------------------------------
	// AI BASE
	// ----------------------------------------------------------------------------------------------
	aiApi := server.RouteGroup(server.ContextUrl("/aibase"))
	{
		aiApi.GET(("/"), server.AddRoute(apis.AiBase))
		aiApi.DELETE(("/"), server.AddRoute(apis.DeleteAiBase))
	}
	// ----------------------------------------------------------------------------------------------
	// Plugin
	// ----------------------------------------------------------------------------------------------
	pluginApi := server.RouteGroup(server.ContextUrl("/plugin"))
	{
		pluginApi.POST(("/service"), server.AddRoute(apis.PluginService))
		pluginApi.GET(("/detail"), server.AddRoute(apis.PluginDetail))
	}

	//
	// 分组管理
	//
	groupApi := server.RouteGroup(server.ContextUrl("/group"))
	{
		groupApi.POST("/create", server.AddRoute(apis.CreateGroup))
		groupApi.PUT("/update", server.AddRoute(apis.UpdateGroup))
		groupApi.GET("/list", server.AddRoute(apis.ListGroup))
		groupApi.GET("/detail", server.AddRoute(apis.GroupDetail))
		groupApi.POST("/bind", server.AddRoute(apis.BindResource))
		groupApi.PUT("/unbind", server.AddRoute(apis.UnBindResource))
		groupApi.DELETE("/", server.AddRoute(apis.DeleteGroup))
	}

	//
	// 协议应用管理
	//
	protoAppApi := server.RouteGroup(server.ContextUrl("/protoapp"))
	{
		protoAppApi.POST("/create", server.AddRoute(apis.CreateProtocolApp))
		protoAppApi.DELETE("/delete", server.AddRoute(apis.DeleteProtocolApp))
		protoAppApi.PUT("/update", server.AddRoute(apis.UpdateProtocolApp))
		protoAppApi.GET("/list", server.AddRoute(apis.ListProtocolApp))
	}
	//
	// 大屏应用管理
	//
	visualApi := server.RouteGroup(server.ContextUrl("/visual"))
	{
		visualApi.POST("/create", server.AddRoute(apis.CreateVisual))
		visualApi.PUT("/publish", server.AddRoute(apis.PublishVisual))
		visualApi.PUT("/update", server.AddRoute(apis.UpdateVisual))
		visualApi.GET("/listByGroup", server.AddRoute(apis.ListVisualByGroup))
		visualApi.GET("/detail", server.AddRoute(apis.VisualDetail))
		visualApi.GET("/group", server.AddRoute(apis.ListVisualGroup))
		visualApi.DELETE("/", server.AddRoute(apis.DeleteVisual))
		// 缩略图
		visualApi.POST("/thumbnail", server.AddRoute(apis.UploadFile))
		visualApi.GET("/thumbnail", server.AddRoute(apis.GetThumbnail))
	}
	/*
	*
	* 模型管理
	*
	 */
	schemaApi := server.RouteGroup(server.ContextUrl("/schema"))
	{
		schemaApi.POST("/create", server.AddRoute(apis.CreateDataSchema))
		schemaApi.DELETE("/delete", server.AddRoute(apis.DeleteDataSchema))
		schemaApi.PUT("/update", server.AddRoute(apis.UpdateDataSchema))
		schemaApi.GET("/list", server.AddRoute(apis.ListDataSchema))
		schemaApi.GET(("/detail"), server.AddRoute(apis.DataSchemaDetail))

	}
	siteConfigApi := server.RouteGroup(server.ContextUrl("/site"))
	{

		siteConfigApi.PUT("/update", server.AddRoute(apis.UpdateSiteConfig))
		siteConfigApi.GET("/detail", server.AddRoute(apis.GetSiteConfig))
	}
	trailerApi := server.RouteGroup(server.ContextUrl("/goods"))
	{
		trailerApi.GET("/list", server.AddRoute(apis.GoodsList))
		trailerApi.GET(("/detail"), server.AddRoute(apis.GoodsDetail))
		trailerApi.POST("/create", server.AddRoute(apis.CreateGoods))
		trailerApi.PUT("/update", server.AddRoute(apis.UpdateGoods))
		trailerApi.POST("/upload", server.AddRoute(apis.UploadGoodsFile))
		trailerApi.DELETE("/delete", server.AddRoute(apis.DeleteGoods))
	}
	dataCenterApi := server.RouteGroup(server.ContextUrl("/dataCenter"))
	{
		dataCenterApi.GET("/schema/detail", server.AddRoute(apis.GetSchemaDetail))
		dataCenterApi.GET("/schema/list", server.AddRoute(apis.GetSchemaList))
		dataCenterApi.GET("/schema/defineList", server.AddRoute(apis.GetSchemaDefineList))
		dataCenterApi.POST("/data/query", server.AddRoute(apis.GetQueryData))
	}
	//
	// 系统设置
	//
	apis.LoadSystemSettingsAPI()

	/**
	 * 定时任务
	 */
	route := server.DefaultApiServer.Route()
	route.StaticFS("api/cron_assets", http.Dir(cron_task.CRON_ASSETS))
	route.StaticFS("api/cron_logs", http.Dir(cron_task.CRON_LOGS))
	crontaskApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/crontask"))
	{
		crontaskApi.POST("/create", server.AddRouteV2(apis.CreateScheduleTask))
		crontaskApi.DELETE("/delete", server.AddRouteV2(apis.DeleteScheduleTask))
		crontaskApi.PUT("/update", server.AddRouteV2(apis.UpdateScheduleTask))
		crontaskApi.GET("/page", server.AddRouteV2(apis.PageScheduleTask))
		crontaskApi.GET("/results/page", server.AddRouteV2(apis.PageCronTaskResult))

		crontaskApi.GET("/start", server.AddRouteV2(apis.EnableTask))
		crontaskApi.GET("/stop", server.AddRouteV2(apis.DisableTask))
		crontaskApi.GET("/listRunningTask", server.AddRouteV2(apis.ListRunningTask))
		crontaskApi.GET("/terminateRunningTask", server.AddRouteV2(apis.TerminateRunningTask))
	}
}

// ApiServerPlugin Start
func (hs *ApiServerPlugin) Start(r typex.RuleX) error {
	hs.ruleEngine = r
	hs.LoadRoute()
	glogger.GLogger.Infof("Http server started on :%v", hs.mainConfig.DbPath)
	return nil
}

func (hs *ApiServerPlugin) Stop() error {
	return nil
}

func (hs *ApiServerPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
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
func (*ApiServerPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{Out: "ApiServerPlugin"}
}
