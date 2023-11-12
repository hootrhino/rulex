package test

import (
	"context"
	"runtime"

	"testing"
	"time"

	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* 测试Trailer的时候比较麻烦，首先建议编译好测试代码
* 建议试试这个脚本: https://github.com/hootrhino/trailer-demo-app.git
*
 */
//  go test -timeout 30s -run ^Test_Trailer_load github.com/hootrhino/rulex/test -v -count=1
//
func Test_Trailer_load(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	path := "script/_temp/trailer-demo-app/trailer-demo-app"
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	goods := trailer.GoodsInfo{
		UUID:        "trailer-demo-app",
		LocalPath:   path,
		NetAddr:     "127.0.0.1:7798",
		Description: "trailer-demo-app",
		Args:        "arg1 args",
	}
	if err := trailer.StartProcess(goods); err != nil {
		glogger.GLogger.Fatal("Goods load failed:", err)
	}
	grpcConnection, err := grpc.Dial(goods.NetAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer grpcConnection.Close()
	client := trailer.NewTrailerClient(grpcConnection)
	schema, err := client.Schema(context.Background(), &trailer.SchemaRequest{})
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	t.Log("=============", schema)
	result, err := client.Query(context.Background(),
		&trailer.DataRowsRequest{Query: []byte("SELECT * FROM TABLE_AAA")})
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	t.Log("=============", result)
	time.Sleep(20 * time.Second)
	engine.Stop()
}
