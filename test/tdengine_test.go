package test

import (
	"context"
	"encoding/json"
	"rulex/core"
	"rulex/engine"
	httpserver "rulex/plugin/http_server"
	"rulex/rulexrpc"
	"rulex/typex"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tdEngineConfig struct {
	Fqdn           string `json:"fqdn" validate:"required"`
	Port           int    `json:"port" validate:"required"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	DbName         string `json:"dbName" validate:"required"`
	CreateDbSql    string `json:"createDbSql" validate:"required"`
	CreateTableSql string `json:"createTableSql" validate:"required"`
	InsertSql      string `json:"insertSql" validate:"required"`
}

func Test_gen_td_config(t *testing.T) {
	td := tdEngineConfig{
		Fqdn:           "127.0.0.1",                                                                      // 服务地址
		Port:           4400,                                                                             // 服务端口
		Username:       "root",                                                                           // 用户
		Password:       "taosdata",                                                                       // 密码
		DbName:         "test",                                                                           // 数据库名
		CreateDbSql:    "CREATE DATABASE IF NOT EXISTS device UPDATE 0;",                                 // 建库SQL
		CreateTableSql: "CREATE TABLE IF NOT EXISTS meter (ts TIMESTAMP, current FLOAT, valtage FLOAT);", // 建表SQL
		InsertSql:      `INSERT INTO meter VALUES (NOW, %v, %v);`,                                        // 插入SQL
	}
	b, _ := json.Marshal(td)
	t.Log(string(b))
}
func Test_gen_tdEngineConfig(t *testing.T) {
	c, err := core.RenderOutConfig(typex.TDENGINE_TARGET, "TDENGINE", tdEngineConfig{})
	if err != nil {
		t.Error(err)
	}
	b, _ := json.MarshalIndent(c.Views, "  ", "")
	t.Log(string(b))
}

/*
*
* Test_data_to_tdengine
*
 */
func Test_data_to_tdengine(t *testing.T) {
	core.InitGlobalConfig()
	engine := engine.NewRuleEngine(core.InitGlobalConfig())
	engine.Start()

	hh := httpserver.NewHttpApiServer(2580, "/../plugin/http_server/www/", "../rulex-test_"+time.Now().Format("2006-01-02-15_04_05")+".db", engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin(hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Test_data_to_tdengine", "Test_data_to_tdengine", map[string]interface{}{
		"port": 2581,
	})

	if err := engine.LoadInEnd(grpcInend); err != nil {
		t.Error("grpcInend load failed:", err)
	}

	tdOutEnd := typex.NewOutEnd(typex.TDENGINE_TARGET,
		"Test_data_to_tdengine",
		"Test_data_to_tdengine",
		map[string]interface{}{
			"fqdn":           "106.15.225.172",
			"port":           6041,
			"username":       "root",
			"password":       "taosdata",
			"dbName":         "device",
			"createDbSql":    "CREATE DATABASE IF NOT EXISTS device UPDATE 0;",
			"createTableSql": "CREATE TABLE IF NOT EXISTS meter (ts TIMESTAMP, current FLOAT, valtage FLOAT);",
			"insertSql":      "INSERT INTO meter VALUES (NOW, %v, %v);",
		})
	if err := engine.LoadOutEnd(tdOutEnd); err != nil {
		t.Error(err)
	}
	//
	// Load Rule
	//
	callback := strings.Replace(
		`Actions = {
			function(data)
				local Result = rulexlib:DataToTdEngineServer('$$UUID', '1024, 1025, 234,33, 30')
				print("rulexlib:DataToTdEngineServer Result", Result==nil)
				return true, data
			end
		}`,
		"$$UUID",
		tdOutEnd.UUID,
		1,
	)
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{grpcInend.UUID},
		`function Success() print("[Test_data_to_tdengine Success Callback]=> OK") end`,
		callback,
		`function Failed(error) print("[Test_data_to_tdengine Failed Callback]", error) end`)

	if err := engine.LoadRule(rule1); err != nil {
		t.Error(err)
	}
	//
	//
	//
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: `[{"co2":10,"hum":30,"lex":22,"temp":100}]`,
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(1 * time.Second)
	engine.Stop()
}
