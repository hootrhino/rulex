package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/rulexrpc"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* Test_data_to_tdengine
*
 */
func Test_data_to_tdengine(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd(
		"GRPC",
		"Test_data_to_tdengine",
		"Test_data_to_tdengine", map[string]interface{}{
			"port": 2581,
			"host": "127.0.0.1",
		})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Fatal("grpcInend load failed:", err)
	}

	tdOutEnd := typex.NewOutEnd(typex.TDENGINE_TARGET,
		"Test_data_to_tdengine",
		"Test_data_to_tdengine",
		map[string]interface{}{
			"fqdn":           "127.0.0.1",
			"port":           6041,
			"username":       "root",
			"password":       "taosdata",
			"dbName":         "device",
			"createDbSql":    "CREATE DATABASE IF NOT EXISTS device UPDATE 0;",
			"createTableSql": "CREATE TABLE IF NOT EXISTS meter01 (ts TIMESTAMP, co2 INT, hum INT, lex INT, temp INT);",
			"insertSql":      "INSERT INTO meter01 VALUES (NOW, %v, %v, %v, %v);",
		})
	tdOutEnd.UUID = "TD1"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(tdOutEnd, ctx1, cancelF1); err != nil {
		t.Fatal(err)
	}
	//
	// Load Rule [{"co2":10,"hum":30,"lex":22,"temp":100}]
	//
	callback :=
		`Actions = {
			function(args)
				local t = rulexlib:J2T(data)
				local Result = data:ToTdEngine('TD1', string.format("%d, %d, %d, %d", t['co2'], t['hum'], t['lex'], t['temp']))
				print("data:ToTdEngine Result", Result==nil)
				return false, data
			end
		}`
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[Test_data_to_tdengine Success Callback]=> OK") end`,
		callback,
		`function Failed(error) print("[Test_data_to_tdengine Failed Callback]", error) end`)

	if err := engine.LoadRule(rule1); err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 3; i++ {
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: fmt.Sprintf(`{"co2":%v,"hum":%v,"lex":%v,"temp":%v}`, rand.Int63n(100), rand.Int63n(100), rand.Int63n(100), rand.Int63n(100)),
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Rulex Rpc Call Result ====>>: %v --%v", resp.GetMessage(), i)

	}

	time.Sleep(3 * time.Second)
	engine.Stop()
}
