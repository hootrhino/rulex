package test

import (
	"rulex/core"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexlib"
	"rulex/sidecar"
	"runtime"

	"testing"
	"time"

	"github.com/ngaut/log"
)

/*
*
* 测试Sidecar的时候比较麻烦，首先建议编译好测试代码
* 建议试试这个脚本: test\script\clone.sh
*
 */
func Test_Sidecar_load(t *testing.T) {
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.StartLogWatcher(core.GlobalConfig.LogPath)
	rulexlib.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "./rulex.db", engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		log.Fatal("Rule load failed:", err)
	}
	path := "script/_temp/grpc_driver_hello_go/grpc_driver_hello_go"
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	if err := engine.LoadGoods(sidecar.Goods{
		UUID:        "grpc_driver_hello_go",
		Addr:        path,
		Description: "grpc_driver_hello_go",
		Args:        []string{"arg1", "arg2"},
	}); err != nil {
		log.Fatal("Goods load failed:", err)
	}

	time.Sleep(5 * time.Second)
	engine.Stop()
}
