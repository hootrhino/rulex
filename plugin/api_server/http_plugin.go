package api_server

import (
	"fmt"
	"github.com/hootrhino/rulex/global"
	"github.com/hootrhino/rulex/plugin/api_server/initialize"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/ini.v1"
	"strconv"
)

var apiPort int = 8000

type HttpPlugin struct {
	uuid string
}

func NewHttpPlugin() *HttpPlugin {
	return &HttpPlugin{
		uuid: "HTTP PLUGIN",
	}
}

func (*HttpPlugin) Init(config *ini.Section) error {
	var apiConfig initialize.HttpConfig
	if err := utils.InIMapToStruct(config, &apiConfig); err != nil {
		return err
	}
	global.RULEX_DB = initialize.Gorm(apiConfig.DbPath)
	apiPort = apiConfig.Port
	if global.RULEX_DB != nil {
		initialize.RegisterTables(global.RULEX_DB) // 初始化表
	}

	return nil
}

func (*HttpPlugin) Start(r typex.RuleX) error {
	router := initialize.Routers(r)
	PORT := strconv.Itoa(apiPort)
	go func() {
		// 启动服务
		if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
			fmt.Println(fmt.Sprintf("服务启动失败:%s", err.Error()))
		}
	}()
	// 此处不能使用优雅退出方式，会导致插件的加载bug
	//exit := make(chan os.Signal)
	//signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	//<-exit
	return nil
}

func (*HttpPlugin) Stop() error {
	return nil
}

func (hp *HttpPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hp.uuid,
		Name:     "RULEX HTTP RESTFul Api Server V2.0.0",
		Version:  "v2.0.0",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "liws",
		Email:    "963755497@qq.com",
		License:  "MIT",
	}
}

func (*HttpPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{Out: "HTTP API SERVER"}
}
