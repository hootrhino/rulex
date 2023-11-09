package test

import (
	"context"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/rulexrpc"
	"github.com/hootrhino/rulex/typex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_rulex_base_lib(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Rulex Grpc InEnd",
		"Rulex Grpc InEnd", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 2581,
		})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Error("grpcInend load failed:", err)
	}
	//
	// Load Rule
	//
	rule := typex.NewRule(engine,
		"uuid4",
		"rule4",
		"rule4",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[Success Callback]=> OK") end`,
		`
	Actions = {
		function(args)
			print("[rulexlib:Time()] ==>", rulexlib:Time())
			print("[rulexlib:TsUnix()] ==>", rulexlib:TsUnix())
			print("[rulexlib:TsUnixNano()] ==>", rulexlib:TsUnixNano())
			local MatchHexS = rulexlib:MatchHex("age:[1,3];sex:[4,5]", "FFFFFF014CB2AA55")
			for key, value in pairs(MatchHexS) do
			    print('rulexlib:MatchHex', key, value)
		    end
			-- rulexlib:VSet('k', 'v')
			-- print("[rulexlib:VGet(k)] ==>", rulexlib:VGet('k'))
			-- Hello()
			-- rulexlib:Throw('this is test Throw')
			return true, args
		end,
		function(args)
			rulexlib:log(rulexlib:Time())
			return true, args
		end
	}`,
		`function Failed(error) print("[Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
	}
	client := rulexrpc.NewRulexRpcClient(conn)

	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: `
				[
					{"co2":10,"hum":30,"lex":22,"temp":100},
					{"co2":100,"hum":300,"lex":220,"temp":1000},
					{"co2":1000,"hum":3000,"lex":2200,"temp":10000}
				]`})

	if err != nil {
		t.Error(err)
	}
	t.Logf("Rulex Rpc Call Result ====>>: %v", resp.GetMessage())

	time.Sleep(5 * time.Second)
	engine.Stop()
}
