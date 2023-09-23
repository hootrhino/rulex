package httpserver

import (
	"encoding/json"
	"net/http"

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
		newGoods := typex.Goods{
			UUID:        mGoods.UUID,
			Addr:        mGoods.Addr,
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
		app := typex.NewApplication(
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
	server.DefaultApiServer.Route().GET(server.ContextUrl("plugins"), server.DefaultApiServer.AddRoute(apis.Plugins))
	//
	// Get system information
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("system"), server.DefaultApiServer.AddRoute(apis.System))
	//
	// Ping -> Pong
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("ping"), server.DefaultApiServer.AddRoute(apis.Ping))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("sourceCount"), server.DefaultApiServer.AddRoute(apis.SourceCount))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("logs"), server.DefaultApiServer.AddRoute(apis.Logs))
	//
	//
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("logout"), server.DefaultApiServer.AddRoute(apis.LogOut))
	//
	// Get all inends
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends"), server.DefaultApiServer.AddRoute(apis.InEnds))
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends/detail"), server.DefaultApiServer.AddRoute(apis.InEndDetail))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("drivers"), server.DefaultApiServer.AddRoute(apis.Drivers))
	//
	// Get all outends
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("outends"), server.DefaultApiServer.AddRoute(apis.OutEnds))
	server.DefaultApiServer.Route().GET(server.ContextUrl("outends/detail"), server.DefaultApiServer.AddRoute(apis.OutEndDetail))
	//
	// Get all rules
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("rules"), server.DefaultApiServer.AddRoute(apis.Rules))
	server.DefaultApiServer.Route().GET(server.ContextUrl("rules/detail"), server.DefaultApiServer.AddRoute(apis.RuleDetail))
	//
	// Get statistics data
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("statistics"), server.DefaultApiServer.AddRoute(apis.Statistics))
	server.DefaultApiServer.Route().GET(server.ContextUrl("snapshot"), server.DefaultApiServer.AddRoute(apis.SnapshotDump))
	//
	// Auth
	//
	userApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/users"))
	{
		userApi.GET(("/"), server.DefaultApiServer.AddRoute(apis.Users))
		userApi.GET(("/detail"), server.DefaultApiServer.AddRoute(apis.UserDetail))
		userApi.POST(("/"), server.DefaultApiServer.AddRoute(apis.CreateUser))
	}

	//
	//
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("login"), server.DefaultApiServer.AddRoute(apis.Login))
	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("info"), server.DefaultApiServer.AddRoute(apis.Info))
	//
	// Create InEnd
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("inends"), server.DefaultApiServer.AddRoute(apis.CreateInend))
	//
	// Update Inend
	//
	server.DefaultApiServer.Route().PUT(server.ContextUrl("inends"), server.DefaultApiServer.AddRoute(apis.UpdateInend))
	//
	// 配置表
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends/config"), server.DefaultApiServer.AddRoute(apis.GetInEndConfig))
	//
	// 数据模型表
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("inends/models"), server.DefaultApiServer.AddRoute(apis.GetInEndModels))
	//
	// Create OutEnd
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("outends"), server.DefaultApiServer.AddRoute(apis.CreateOutEnd))
	//
	// Update OutEnd
	//
	server.DefaultApiServer.Route().PUT(server.ContextUrl("outends"), server.DefaultApiServer.AddRoute(apis.UpdateOutEnd))
	rulesApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/rules"))
	{
		rulesApi.POST(("/"), server.DefaultApiServer.AddRoute(apis.CreateRule))
		rulesApi.PUT(("/"), server.DefaultApiServer.AddRoute(apis.UpdateRule))
		rulesApi.DELETE(("/"), server.DefaultApiServer.AddRoute(apis.DeleteRule))
		rulesApi.POST(("/testIn"), server.DefaultApiServer.AddRoute(apis.TestSourceCallback))
		rulesApi.POST(("/testOut"), server.DefaultApiServer.AddRoute(apis.TestOutEndCallback))
		rulesApi.POST(("/testDevice"), server.DefaultApiServer.AddRoute(apis.TestDeviceCallback))
	}

	//
	// Delete inend by UUID
	//
	server.DefaultApiServer.Route().DELETE(server.ContextUrl("inends"), server.DefaultApiServer.AddRoute(apis.DeleteInEnd))
	//
	// Delete outEnd by UUID
	//
	server.DefaultApiServer.Route().DELETE(server.ContextUrl("outends"), server.DefaultApiServer.AddRoute(apis.DeleteOutEnd))

	//
	// 验证 lua 语法
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("validateRule"), server.DefaultApiServer.AddRoute(apis.ValidateLuaSyntax))
	//
	// 获取配置表
	//
	resourceTypeApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/"))
	{
		resourceTypeApi.GET(("rType"), server.DefaultApiServer.AddRoute(apis.RType))
		resourceTypeApi.GET(("tType"), server.DefaultApiServer.AddRoute(apis.TType))
		resourceTypeApi.GET(("dType"), server.DefaultApiServer.AddRoute(apis.DType))
	}
	//
	// 网络适配器列表
	//
	osApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/os"))
	{
		osApi.GET(("/netInterfaces"), server.DefaultApiServer.AddRoute(apis.GetNetInterfaces))
		osApi.GET(("/osRelease"), server.DefaultApiServer.AddRoute(apis.CatOsRelease))
		osApi.GET(("/uarts"), server.DefaultApiServer.AddRoute(apis.GetUartList))
		osApi.GET(("/system"), server.DefaultApiServer.AddRoute(apis.System))
		osApi.GET(("/startedAt"), server.DefaultApiServer.AddRoute(apis.StartedAt))

	}
	//
	// 设备管理
	//
	deviceApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/devices"))
	{
		deviceApi.GET(("/"), server.DefaultApiServer.AddRoute(apis.Devices))
		deviceApi.POST(("/"), server.DefaultApiServer.AddRoute(apis.CreateDevice))
		deviceApi.PUT(("/"), server.DefaultApiServer.AddRoute(apis.UpdateDevice))
		deviceApi.DELETE(("/"), server.DefaultApiServer.AddRoute(apis.DeleteDevice))
		deviceApi.GET(("/detail"), server.DefaultApiServer.AddRoute(apis.DeviceDetail))
		deviceApi.POST(("/modbus/sheetImport"), server.DefaultApiServer.AddRoute(apis.ModbusSheetImport))
		deviceApi.PUT(("/modbus/point"), server.DefaultApiServer.AddRoute(apis.UpdateModbusPoint))
		deviceApi.GET(("/modbus"), server.DefaultApiServer.AddRoute(apis.ModbusPoints))
	}
	goodsApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/goods"))
	{
		// 外挂管理
		goodsApi.GET(("/"), server.DefaultApiServer.AddRoute(apis.Goods))
		goodsApi.POST(("/"), server.DefaultApiServer.AddRoute(apis.CreateGoods))
		goodsApi.PUT(("/"), server.DefaultApiServer.AddRoute(apis.UpdateGoods))
		goodsApi.DELETE(("/"), server.DefaultApiServer.AddRoute(apis.DeleteGoods))
	}

	// ----------------------------------------------------------------------------------------------
	// APP
	// ----------------------------------------------------------------------------------------------
	appApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/app"))
	{
		appApi.GET(("/"), server.DefaultApiServer.AddRoute(apis.Apps))
		appApi.POST(("/"), server.DefaultApiServer.AddRoute(apis.CreateApp))
		appApi.PUT(("/"), server.DefaultApiServer.AddRoute(apis.UpdateApp))
		appApi.DELETE(("/"), server.DefaultApiServer.AddRoute(apis.RemoveApp))
		appApi.PUT(("/start"), server.DefaultApiServer.AddRoute(apis.StartApp))
		appApi.PUT(("/stop"), server.DefaultApiServer.AddRoute(apis.StopApp))
		appApi.GET(("/detail"), server.DefaultApiServer.AddRoute(apis.AppDetail))
	}
	// ----------------------------------------------------------------------------------------------
	// AI BASE
	// ----------------------------------------------------------------------------------------------
	aiApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/aibase"))
	{
		aiApi.GET(("/"), server.DefaultApiServer.AddRoute(apis.AiBase))
		aiApi.DELETE(("/"), server.DefaultApiServer.AddRoute(apis.DeleteAiBase))
	}
	// ----------------------------------------------------------------------------------------------
	// Plugin
	// ----------------------------------------------------------------------------------------------
	pluginApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/plugin"))
	{
		pluginApi.POST(("/service"), server.DefaultApiServer.AddRoute(apis.PluginService))
		pluginApi.GET(("/detail"), server.DefaultApiServer.AddRoute(apis.PluginDetail))
	}

	//
	// 分组管理
	//
	groupApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/group"))
	{
		groupApi.POST("/create", server.DefaultApiServer.AddRoute(apis.CreateGroup))
		groupApi.DELETE("/delete", server.DefaultApiServer.AddRoute(apis.DeleteGroup))
		groupApi.PUT("/update", server.DefaultApiServer.AddRoute(apis.UpdateGroup))
		groupApi.GET("/list", server.DefaultApiServer.AddRoute(apis.ListGroup))
		groupApi.POST("/bind", server.DefaultApiServer.AddRoute(apis.BindResource))
		groupApi.PUT("/unbind", server.DefaultApiServer.AddRoute(apis.UnBindResource))
		groupApi.GET("/devices", server.DefaultApiServer.AddRoute(apis.FindDeviceByGroup))
		groupApi.GET("/visuals", server.DefaultApiServer.AddRoute(apis.FindVisualByGroup))
	}

	//
	// 协议应用管理
	//
	protoAppApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/protoapp"))
	{
		protoAppApi.POST("/create", server.DefaultApiServer.AddRoute(apis.CreateProtocolApp))
		protoAppApi.DELETE("/delete", server.DefaultApiServer.AddRoute(apis.DeleteProtocolApp))
		protoAppApi.PUT("/update", server.DefaultApiServer.AddRoute(apis.UpdateProtocolApp))
		protoAppApi.GET("/list", server.DefaultApiServer.AddRoute(apis.ListProtocolApp))
	}
	//
	// 大屏应用管理
	//
	screenApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/visual"))
	{
		screenApi.POST("/create", server.DefaultApiServer.AddRoute(apis.CreateVisual))
		screenApi.DELETE("/delete", server.DefaultApiServer.AddRoute(apis.DeleteVisual))
		screenApi.PUT("/update", server.DefaultApiServer.AddRoute(apis.UpdateVisual))
		screenApi.GET("/list", server.DefaultApiServer.AddRoute(apis.ListVisual))
		screenApi.GET("/GenComponentUUID", server.DefaultApiServer.AddRoute(apis.GenComponentUUID))
	}
	/*
	*
	* 模型管理
	*
	 */
	schemaApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/schema"))
	{
		schemaApi.POST("/create", server.DefaultApiServer.AddRoute(apis.CreateDataSchema))
		schemaApi.DELETE("/delete", server.DefaultApiServer.AddRoute(apis.DeleteDataSchema))
		schemaApi.PUT("/update", server.DefaultApiServer.AddRoute(apis.UpdateDataSchema))
		schemaApi.GET("/list", server.DefaultApiServer.AddRoute(apis.ListDataSchema))
		schemaApi.GET(("/detail"), server.DefaultApiServer.AddRoute(apis.DataSchemaDetail))

	}
	siteConfigApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/site"))
	{

		siteConfigApi.PUT("/update", server.DefaultApiServer.AddRoute(apis.UpdateSiteConfig))
		siteConfigApi.GET("/detail", server.DefaultApiServer.AddRoute(apis.GetSiteConfig))
	}
	//
	// 系统设置
	//
	apis.LoadSystemSettingsAPI()

	/**
	 * 定时任务
	 */
	route := server.DefaultApiServer.Route()
	route.StaticFS("api/cron_assets", http.Dir("cron_asserts"))
	route.StaticFS("api/cron_logs", http.Dir("cron_logs"))
	crontaskApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/crontask"))
	{
		crontaskApi.POST("/create", server.DefaultApiServer.AddRouteV2(apis.CreateScheduleTask))
		crontaskApi.DELETE("/delete", server.DefaultApiServer.AddRouteV2(apis.DeleteScheduleTask))
		crontaskApi.PUT("/update", server.DefaultApiServer.AddRouteV2(apis.UpdateScheduleTask))
		crontaskApi.GET("/page", server.DefaultApiServer.AddRouteV2(apis.PageScheduleTask))

		crontaskApi.GET("/start", server.DefaultApiServer.AddRouteV2(apis.EnableTask))
		crontaskApi.GET("/stop", server.DefaultApiServer.AddRouteV2(apis.DisableTask))
		crontaskApi.GET("/listRunningTask", server.DefaultApiServer.AddRouteV2(apis.ListRunningTask))
		crontaskApi.GET("/terminateRunningTask", server.DefaultApiServer.AddRouteV2(apis.TerminateRunningTask))

		crontaskApi.GET("/results/page", server.DefaultApiServer.AddRouteV2(apis.PageCronTaskResult))
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
