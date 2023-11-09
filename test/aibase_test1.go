package test

import (
	"context"

	"github.com/hootrhino/rulex/component/rulexrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"testing"
	"time"

	"github.com/hootrhino/rulex/typex"
)

// go test -timeout 30s -run ^Test_AIBASE_ANN_MNIST github.com/hootrhino/rulex/test -v -count=1

func Test_AIBASE_ANN_MNIST(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "Test_AIBASE_ANN_MNIST",
		"Test_AIBASE_ANN_MNIST", map[string]interface{}{
			"host": "127.0.0.1",
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
			    local P = {
				    [0] = {11,12,13,14,15,16,17,18},
				    [1] = {21,22,23,24,25,26,27,28},
				    [2] = {31,32,33,34,35,36,37,38}
				}
				local result, err1 = aibase:Infer('BUILDIN_MNIST', P)
				for index, value in ipairs(result) do
					for index2, value2 in ipairs(value) do
						print(index, index2, value2)
					end
				end
				print('Test_AIBASE_ANN_MNIST =>', result, err1)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)
	for i := 0; i < 2; i++ {
		resp, err := client.Work(context.Background(), &rulexrpc.Data{
			Value: `{"value":"0298010d"}`,
		})
		if err != nil {
			t.Error(err)
		}
		t.Logf("Rulex Rpc Call Result ====>>: %v --%v", resp.GetMessage(), i)

	}

	time.Sleep(5 * time.Second)
	engine.Stop()
}
