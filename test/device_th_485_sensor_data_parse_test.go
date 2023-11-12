package test

import (
	"context"
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
* Test_485_sensor
*
 */
func Test_modbus_485_sensor_data_parse(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Test_485_sensor", "Test_485_sensor", map[string]interface{}{
		"port": 2581,
	})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Error("grpcInend load failed:", err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
				local table = rulexlib:J2T(data)
				local value = table['value']
				local t = rulexlib:HsubToN(value, 5, 8)
				local h = rulexlib:HsubToN(value, 0, 4)
				local t1 = rulexlib:HToN(string.sub(value, 5, 8))
				local h2 = rulexlib:HToN(string.sub(value, 0, 4))
				print('Data ========> ', rulexlib:T2J({
					Device = "TH00000001",
					Ts = rulexlib:TsUnix(),
					T = t,
					H = h,
					T1 = t1,
					H2 = h2
				}))
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2; i++ {
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: `
			{
				"tag":"data",
				"function":3,
				"address":0,
				"quantity":4,
				"value":"0298010d"
			}
			`,
		})
		if err != nil {
			t.Error(err)
		}
		t.Logf("Rulex Rpc Call Result ====>>: %v --%v", resp.GetMessage(), i)

	}

	time.Sleep(3 * time.Second)
	engine.Stop()
}
