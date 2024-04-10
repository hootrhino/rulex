package httpserver

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hootrhino/rulex/component/cron_task"
	"github.com/hootrhino/rulex/component/hwportmanager"
	"github.com/shirou/gopsutil/cpu"

	"github.com/hootrhino/rulex/component/appstack"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/rulex_api_server/apis"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/server"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/component/trailer"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"gopkg.in/ini.v1"

	_ "github.com/mattn/go-sqlite3"
)

type _serverConfig struct {
	DbPath string `ini:"dbpath"`
	Port   int    `ini:"port"`
}
type ApiServerPlugin struct {
	uuid       string
	ruleEngine typex.RuleX
	mainConfig _serverConfig
}

func NewHttpApiServer(ruleEngine typex.RuleX) *ApiServerPlugin {
	return &ApiServerPlugin{
		uuid:       "HTTP-API-SERVER",
		mainConfig: _serverConfig{Port: 2580},
		ruleEngine: ruleEngine,
	}
}

/*
*
* 初始化RULEX, 初始化数据到运行时
*
 */
func initRulex(engine typex.RuleX) {
	go GetCpuUsage()
	/*
	*
	* 加载Port
	*
	 */
	loadAllPortConfig()
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
		newGoods := trailer.GoodsInfo{
			UUID:        mGoods.UUID,
			AutoStart:   mGoods.AutoStart,
			LocalPath:   mGoods.LocalPath,
			NetAddr:     mGoods.NetAddr,
			Args:        mGoods.Args,
			ExecuteType: mGoods.ExecuteType,
			Description: mGoods.Description,
		}
		if err := trailer.StartProcess(newGoods); err != nil {
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
			glogger.GLogger.Debug("App autoStart allowed:", app.UUID)
			if err1 := appstack.StartApp(app.UUID); err1 != nil {
				glogger.GLogger.Error("App autoStart failed:", err1)
			}
		}
	}
	//
	// load Cron Task
	for _, task := range service.AllEnabledCronTask() {
		if err := cron_task.GetCronManager().AddTask(task); err != nil {
			glogger.GLogger.Error(err)
			continue
		}
	}

}

/*
*
* 从数据库拿端口配置
*
 */
func loadAllPortConfig() {
	MHwPorts, err := service.AllHwPort()
	if err != nil {
		glogger.GLogger.Fatal(err)
		return
	}
	for _, MHwPort := range MHwPorts {
		Port := hwportmanager.RhinoH3HwPort{
			UUID:        MHwPort.UUID,
			Name:        MHwPort.Name,
			Type:        MHwPort.Type,
			Alias:       MHwPort.Alias,
			Description: MHwPort.Description,
		}
		// 串口
		if MHwPort.Type == "UART" {
			config := hwportmanager.UartConfig{}
			if err := utils.BindConfig(MHwPort.GetConfig(), &config); err != nil {
				glogger.GLogger.Error(err) // 这里必须不能出错
				continue
			}
			Port.Config = config
			hwportmanager.SetHwPort(Port)
		}
		// 未知接口参数为空，以后扩展，比如FD
		if MHwPort.Type != "UART" {
			Port.Config = nil
			hwportmanager.SetHwPort(Port)
		}
	}
}

func (hs *ApiServerPlugin) Init(config *ini.Section) error {
	if err := utils.InIMapToStruct(config, &hs.mainConfig); err != nil {
		return err
	}
	server.StartRulexApiServer(hs.ruleEngine, hs.mainConfig.Port)

	interdb.DB().Exec("VACUUM;")
	interdb.RegisterModel(
		&model.MInEnd{},
		&model.MOutEnd{},
		&model.MRule{},
		&model.MUser{},
		&model.MDevice{},
		&model.MGoods{},
		&model.MApp{},
		&model.MVisual{},
		&model.MGenericGroup{},
		&model.MGenericGroupRelation{},
		&model.MNetworkConfig{},
		&model.MWifiConfig{},
		&model.MIotSchema{},
		&model.MIotProperty{},
		&model.MIpRoute{},
		&model.MCronTask{},
		&model.MCronResult{},
		&model.MHwPort{},
		&model.MInternalNotify{},
		&model.MUserLuaTemplate{},
		&model.MModbusDataPoint{},
		&model.MSiemensDataPoint{},
		&model.MHnc8DataPoint{},
		&model.MKnd8DataPoint{},
	)
	// 初始化所有预制参数
	server.DefaultApiServer.InitializeGenericOSData()
	server.DefaultApiServer.InitializeEEKITData()
	server.DefaultApiServer.InitializeWindowsData()
	server.DefaultApiServer.InitializeUnixData()
	server.DefaultApiServer.InitializeConfigCtl()
	initRulex(hs.ruleEngine)
	return nil
}

/*
*
* 加载路由
*
 */
func (hs *ApiServerPlugin) LoadRoute() {
	systemApi := server.RouteGroup(server.ContextUrl("/"))
	{
		systemApi.GET(("/ping"), server.AddRoute(apis.Ping))
	}

	//
	//
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("drivers"), server.AddRoute(apis.Drivers))

	//
	// Get statistics data
	//
	server.DefaultApiServer.Route().GET(server.ContextUrl("statistics"), server.AddRoute(apis.Statistics))
	//
	// Auth
	//
	userApi := server.RouteGroup(server.ContextUrl("/users"))
	{
		// userApi.GET(("/"), server.AddRoute(apis.Users))
		userApi.POST(("/"), server.AddRoute(apis.CreateUser))
		userApi.PUT(("/update"), server.AddRoute(apis.UpdateUser))
		userApi.GET(("/detail"), server.AddRoute(apis.UserDetail))
		userApi.POST(("/logout"), server.AddRoute(apis.LogOut))

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
	InEndApi := server.RouteGroup(server.ContextUrl("/inends"))
	{
		InEndApi.GET(("/detail"), server.AddRoute(apis.InEndDetail))
		InEndApi.GET(("/list"), server.AddRoute(apis.InEnds))
		InEndApi.POST(("/create"), server.AddRoute(apis.CreateInend))
		InEndApi.DELETE(("/del"), server.AddRoute(apis.DeleteInEnd))
		InEndApi.PUT(("/update"), server.AddRoute(apis.UpdateInend))
		InEndApi.PUT("/restart", server.AddRoute(apis.RestartInEnd))
	}

	rulesApi := server.RouteGroup(server.ContextUrl("/rules"))
	{
		rulesApi.POST(("/create"), server.AddRoute(apis.CreateRule))
		rulesApi.PUT(("/update"), server.AddRoute(apis.UpdateRule))
		rulesApi.DELETE(("/del"), server.AddRoute(apis.DeleteRule))
		rulesApi.GET(("/list"), server.AddRoute(apis.Rules))
		rulesApi.GET(("/detail"), server.AddRoute(apis.RuleDetail))
		//
		rulesApi.POST(("/testIn"), server.AddRoute(apis.TestSourceCallback))
		rulesApi.POST(("/testOut"), server.AddRoute(apis.TestOutEndCallback))
		rulesApi.POST(("/testDevice"), server.AddRoute(apis.TestDeviceCallback))
		rulesApi.GET(("/byInend"), server.AddRoute(apis.ListByInend))
		rulesApi.GET(("/byDevice"), server.AddRoute(apis.ListByDevice))
		//
		rulesApi.GET(("/getCanUsedResources"), server.AddRoute(apis.GetAllResources))
		//
		rulesApi.POST(("/formatLua"), server.AddRoute(apis.FormatLua))

	}
	OutEndApi := server.RouteGroup(server.ContextUrl("/outends"))
	{
		OutEndApi.GET(("/detail"), server.AddRoute(apis.OutEndDetail))
		OutEndApi.GET(("/list"), server.AddRoute(apis.OutEnds))
		OutEndApi.POST(("/create"), server.AddRoute(apis.CreateOutEnd))
		OutEndApi.DELETE(("/del"), server.AddRoute(apis.DeleteOutEnd))
		OutEndApi.PUT(("/update"), server.AddRoute(apis.UpdateOutEnd))
		OutEndApi.PUT("/restart", server.AddRoute(apis.RestartOutEnd))
	}

	//
	// 验证 lua 语法
	//
	server.DefaultApiServer.Route().POST(server.ContextUrl("validateRule"), server.AddRoute(apis.ValidateLuaSyntax))

	//
	// 网络适配器列表
	//
	osApi := server.RouteGroup(server.ContextUrl("/os"))
	{
		osApi.GET(("/netInterfaces"), server.AddRoute(apis.GetNetInterfaces))
		osApi.GET(("/osRelease"), server.AddRoute(apis.CatOsRelease))
		osApi.GET(("/system"), server.AddRoute(apis.System))
		osApi.GET(("/startedAt"), server.AddRoute(apis.StartedAt))
		osApi.GET(("/getVideos"), server.AddRoute(apis.GetVideos))
		osApi.GET(("/getGpuInfo"), server.AddRoute(apis.GetGpuInfo))
		osApi.POST(("/resetInterMetric"), server.AddRoute(apis.ResetInterMetric))
	}
	backupApi := server.RouteGroup(server.ContextUrl("/backup"))
	{
		backupApi.GET(("/download"), server.AddRoute(apis.DownloadSqlite))
		backupApi.POST(("/upload"), server.AddRoute(apis.UploadSqlite))
		backupApi.GET(("/snapshot"), server.AddRoute(apis.SnapshotDump))
		backupApi.GET(("/runningLog"), server.AddRoute(apis.GetRunningLog))
	}
	//
	// 设备管理
	//
	deviceApi := server.RouteGroup(server.ContextUrl("/devices"))
	{
		deviceApi.POST(("/create"), server.AddRoute(apis.CreateDevice))
		deviceApi.PUT(("/update"), server.AddRoute(apis.UpdateDevice))
		deviceApi.DELETE(("/del"), server.AddRoute(apis.DeleteDevice))
		deviceApi.GET(("/detail"), server.AddRoute(apis.DeviceDetail))
		deviceApi.GET("/group", server.AddRoute(apis.ListDeviceGroup))
		deviceApi.GET("/listByGroup", server.AddRoute(apis.ListDeviceByGroup))
		deviceApi.PUT("/restart", server.AddRoute(apis.RestartDevice))
		deviceApi.GET("/properties", server.AddRoute(apis.DevicePropertiesPage))
		deviceApi.GET("/deviceErrMsg", server.AddRoute(apis.GetDeviceErrorMsg))
		deviceApi.GET("/pointErrMsg", server.AddRoute(apis.GetDevicePointErrorMsg))
	}
	// Modbus 点位表
	modbusApi := server.RouteGroup(server.ContextUrl("/modbus_data_sheet"))
	{
		modbusApi.POST(("/sheetImport"), server.AddRoute(apis.ModbusSheetImport))
		modbusApi.GET(("/sheetExport"), server.AddRoute(apis.ModbusPointsExport))
		modbusApi.GET(("/list"), server.AddRoute(apis.ModbusSheetPageList))
		modbusApi.POST(("/update"), server.AddRoute(apis.ModbusSheetUpdate))
		modbusApi.DELETE(("/delIds"), server.AddRoute(apis.ModbusSheetDelete))
		modbusApi.DELETE(("/delAll"), server.AddRoute(apis.ModbusSheetDeleteAll))
	}
	// S1200 点位表
	SIEMENS_PLC := server.RouteGroup(server.ContextUrl("/s1200_data_sheet"))
	{
		SIEMENS_PLC.POST(("/sheetImport"), server.AddRoute(apis.SiemensSheetImport))
		SIEMENS_PLC.GET(("/sheetExport"), server.AddRoute(apis.SiemensPointsExport))
		SIEMENS_PLC.GET(("/list"), server.AddRoute(apis.SiemensSheetPageList))
		SIEMENS_PLC.POST(("/update"), server.AddRoute(apis.SiemensSheetUpdate))
		SIEMENS_PLC.DELETE(("/delIds"), server.AddRoute(apis.SiemensSheetDelete))
		SIEMENS_PLC.DELETE(("/delAll"), server.AddRoute(apis.SiemensSheetDeleteAll))
	}
	// 华中数控 点位表
	Hnc8 := server.RouteGroup(server.ContextUrl("/hnc8_data_sheet"))
	{
		Hnc8.POST(("/sheetImport"), server.AddRoute(apis.Hnc8SheetImport))
		Hnc8.GET(("/sheetExport"), server.AddRoute(apis.Hnc8PointsExport))
		Hnc8.GET(("/list"), server.AddRoute(apis.Hnc8SheetPageList))
		Hnc8.POST(("/update"), server.AddRoute(apis.Hnc8SheetUpdate))
		Hnc8.DELETE(("/delIds"), server.AddRoute(apis.Hnc8SheetDelete))
		Hnc8.DELETE(("/delAll"), server.AddRoute(apis.Hnc8SheetDeleteAll))
	}

	// ----------------------------------------------------------------------------------------------
	// APP
	// ----------------------------------------------------------------------------------------------
	appApi := server.RouteGroup(server.ContextUrl("/app"))
	{
		appApi.GET(("/list"), server.AddRoute(apis.Apps))
		appApi.POST(("/create"), server.AddRoute(apis.CreateApp))
		appApi.PUT(("/update"), server.AddRoute(apis.UpdateApp))
		appApi.DELETE(("/del"), server.AddRoute(apis.RemoveApp))
		appApi.PUT(("/start"), server.AddRoute(apis.StartApp))
		appApi.PUT(("/stop"), server.AddRoute(apis.StopApp))
		appApi.GET(("/detail"), server.AddRoute(apis.AppDetail))
	}
	// ----------------------------------------------------------------------------------------------
	// Plugin
	// ----------------------------------------------------------------------------------------------
	pluginsApi := server.RouteGroup(server.ContextUrl("/plugware"))
	{
		pluginsApi.GET(("/list"), server.AddRoute(apis.Plugins))
		pluginsApi.POST(("/service"), server.AddRoute(apis.PluginService))
		pluginsApi.GET(("/detail"), server.AddRoute(apis.PluginDetail))
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
		groupApi.DELETE("/del", server.AddRoute(apis.DeleteGroup))
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
	//
	// 用户LUA代码段管理
	//
	userLuaApi := server.RouteGroup(server.ContextUrl("/userlua"))
	{
		userLuaApi.POST("/create", server.AddRoute(apis.CreateUserLuaTemplate))
		userLuaApi.PUT("/update", server.AddRoute(apis.UpdateUserLuaTemplate))
		userLuaApi.GET("/listByGroup", server.AddRoute(apis.ListUserLuaTemplateByGroup))
		userLuaApi.GET("/detail", server.AddRoute(apis.UserLuaTemplateDetail))
		userLuaApi.GET("/group", server.AddRoute(apis.ListUserLuaTemplateGroup))
		userLuaApi.DELETE("/del", server.AddRoute(apis.DeleteUserLuaTemplate))
		userLuaApi.GET("/search", server.AddRoute(apis.SearchUserLuaTemplateGroup))
	}
	/*
	*
	* 模型管理
	*
	 */
	schemaApi := server.RouteGroup(server.ContextUrl("/schema"))
	{
		// 物模型
		schemaApi.POST("/create", server.AddRoute(apis.CreateDataSchema))
		schemaApi.DELETE("/del", server.AddRoute(apis.DeleteDataSchema))
		schemaApi.PUT("/update", server.AddRoute(apis.UpdateDataSchema))
		schemaApi.GET("/list", server.AddRoute(apis.ListDataSchema))
		schemaApi.GET(("/detail"), server.AddRoute(apis.DataSchemaDetail))
		// 属性
		schemaApi.POST(("/properties/create"), server.AddRoute(apis.CreateIotSchemaProperty))
		schemaApi.PUT(("/properties/update"), server.AddRoute(apis.UpdateIotSchemaProperty))
		schemaApi.DELETE(("/properties/del"), server.AddRoute(apis.DeleteIotSchemaProperty))
		schemaApi.GET(("/properties/list"), server.AddRoute(apis.IotSchemaPropertyPageList))
		schemaApi.GET(("/properties/detail"), server.AddRoute(apis.IotSchemaPropertyDetail))

	}
	trailerApi := server.RouteGroup(server.ContextUrl("/goods"))
	{
		trailerApi.GET("/list", server.AddRoute(apis.GoodsList))
		trailerApi.GET(("/detail"), server.AddRoute(apis.GoodsDetail))
		trailerApi.POST("/create", server.AddRoute(apis.CreateGoods))
		trailerApi.PUT("/update", server.AddRoute(apis.UpdateGoods))
		trailerApi.PUT("/cleanGarbage", server.AddRoute(apis.CleanGoodsUpload))
		trailerApi.PUT("/start", server.AddRoute(apis.StartGoods))
		trailerApi.PUT("/stop", server.AddRoute(apis.StopGoods))
		trailerApi.DELETE("/", server.AddRoute(apis.DeleteGoods))
	}
	// 数据中心
	dataCenterApi := server.RouteGroup(server.ContextUrl("/dataCenter"))
	{
		dataCenterApi.GET("/schema/define", server.AddRoute(apis.GetSchemaDefine))
		dataCenterApi.GET("/schema/detail", server.AddRoute(apis.GetSchemaDetail))
		dataCenterApi.GET("/schema/list", server.AddRoute(apis.GetSchemaList))
		dataCenterApi.GET("/schema/defineList", server.AddRoute(apis.GetSchemaDefineList))
		dataCenterApi.POST("/data/query", server.AddRoute(apis.GetQueryData))
	}
	// 硬件接口API
	HwIFaceApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/hwiface"))
	{
		HwIFaceApi.GET("/detail", server.AddRoute(apis.GetHwPortDetail))
		HwIFaceApi.GET("/list", server.AddRoute(apis.AllHwPorts))
		HwIFaceApi.POST("/update", server.AddRoute(apis.UpdateHwPortConfig))
		HwIFaceApi.GET("/refresh", server.AddRoute(apis.RefreshPortList))
	}
	// 站内公告
	internalNotifyApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/notify"))
	{
		internalNotifyApi.GET("/header", server.AddRoute(apis.InternalNotifiesHeader))
		internalNotifyApi.GET("/list", server.AddRoute(apis.InternalNotifies))
		internalNotifyApi.PUT("/clear", server.AddRoute(apis.ClearInternalNotifies))
		internalNotifyApi.PUT("/read", server.AddRoute(apis.ReadInternalNotifies))
		// internalNotifyApi.POST("/test", server.AddRoute(apis.TestCreateNotifies))
	}
	//
	// 系统设置
	//
	apis.LoadSystemSettingsAPI()

	/**
	 * 定时任务
	 */
	crontaskApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/crontask"))
	{
		crontaskApi.POST("/create", server.AddRouteV2(apis.CreateCronTask))
		crontaskApi.DELETE("/del", server.AddRouteV2(apis.DeleteCronTask))
		crontaskApi.PUT("/update", server.AddRouteV2(apis.UpdateCronTask))
		crontaskApi.GET("/list", server.AddRouteV2(apis.ListCronTask))
		crontaskApi.GET("/results/page", server.AddRouteV2(apis.PageCronTaskResult))
		crontaskApi.GET("/start", server.AddRouteV2(apis.StartTask))
		crontaskApi.GET("/stop", server.AddRouteV2(apis.StopTask))
	}
	//
	// jpegStream APi
	//
	jpegStream := server.DefaultApiServer.GetGroup(server.ContextUrl("/jpeg_stream"))
	{
		jpegStream.GET("/list", server.AddRoute(apis.GetJpegStreamList))
		jpegStream.GET("/detail", server.AddRoute(apis.GetJpegStreamDetail))
	}

	/**
	  swagger config
	  @reference http://localhost:2580/swagger/index.html
	*/
	route := server.DefaultApiServer.Route()
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

// ApiServerPlugin Start
func (hs *ApiServerPlugin) Start(r typex.RuleX) error {
	hs.ruleEngine = r
	hs.LoadRoute()
	glogger.GLogger.Infof("Http server started on :%v", 2580)
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
		Homepage: "/",
		HelpLink: "/",
		Author:   "RHILEXTeam",
		Email:    "RHILEXTeam@hootrhino.com",
		License:  "AGPL",
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
func GetCpuUsage() {
	for {
		select {
		case <-typex.GCTX.Done():
			{
				return
			}
		default:
			{
			}
		}
		cpuPercent, _ := cpu.Percent(time.Duration(10)*time.Second, true)
		V := calculateCpuPercent(cpuPercent)
		if V > 95 {
			service.InsertInternalNotify(model.MInternalNotify{
				UUID:    utils.MakeUUID("NOTIFY"), // UUID
				Type:    `WARNING`,                // INFO | ERROR | WARNING
				Status:  1,
				Event:   `system.cpu.load`, // 字符串
				Ts:      uint64(time.Now().UnixMilli()),
				Summary: "CPU负载过高",
				Info:    fmt.Sprintf("CPU负载过高: %.2f%%, 请注意维护设备", V),
			})
		}
	}

}

// 计算CPU平均使用率
func calculateCpuPercent(cpus []float64) float64 {
	var acc float64 = 0
	for _, v := range cpus {
		acc += v
	}
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", acc/float64(len(cpus))), 64)
	return value
}
