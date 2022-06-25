package test

import (
	"runtime"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	httpserver "github.com/i4de/rulex/plugin/http_server"
	"github.com/i4de/rulex/sidecar"

	"testing"
	"time"
)

/*
*
* 测试Sidecar的时候比较麻烦，首先建议编译好测试代码
* 建议试试这个脚本: test\script\clone.sh
*
 */
func Test_Sidecar_load(t *testing.T) {
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	hh := httpserver.NewHttpApiServer()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
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
		glogger.GLogger.Fatal("Goods load failed:", err)
	}

	time.Sleep(5 * time.Second)
	engine.Stop()
}
