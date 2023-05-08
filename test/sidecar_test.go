package test

import (
	"runtime"

	"testing"
	"time"

	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 测试Trailer的时候比较麻烦，首先建议编译好测试代码
* 建议试试这个脚本: test\script\clone.sh
*
 */
func Test_Trailer_load(t *testing.T) {
	engine := RunTestEngine()
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
	if err := engine.LoadGoods(typex.Goods{
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
