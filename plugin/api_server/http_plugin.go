package api_server

import (
	"fmt"
	"github.com/hootrhino/rulex/global"
	"github.com/hootrhino/rulex/plugin/api_server/initialize"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/ini.v1"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

func (*HttpPlugin) Start(typex.RuleX) error {
	// 优雅退出程序
	router := initialize.Routers()
	PORT := strconv.Itoa(apiPort)
	go func() {
		// 启动服务
		if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
			fmt.Println(fmt.Sprintf("服务启动失败:%s", err.Error()))
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	return nil
}

func (*HttpPlugin) Stop() error {
	return nil
}

func (hp *HttpPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hp.uuid,
		Name:     "RULEX HTTP RESTFul Api Server",
		Version:  "v2.0.0",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

func (*HttpPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{Out: "HTTP API SERVER"}
}
