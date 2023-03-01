package test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	httpserver "github.com/i4de/rulex/plugin/http_server"
	"github.com/i4de/rulex/rulexrpc"
	"github.com/i4de/rulex/typex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FFFFFF014CB2AA55
// go test -timeout 30s -run ^TestHexEncoding github.com/i4de/rulex/test -v -count=1
func TestHexEncoding(t *testing.T) {
	hexs := []byte{255, 255, 255, 1, 76, 178, 170, 85}
	s := hex.EncodeToString(hexs)
	t.Log(fmt.Sprintf("%X", hexs) == s)
	t.Log(fmt.Sprintf("%x", hexs) == s)
	t.Log(hex.DecodeString(s))
}

// go test -timeout 30s -run ^TestCustomProtocolDevice github.com/i4de/rulex/test -v -count=1

func TestCustomProtocolDevice(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer()); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd(typex.GRPC,
		"Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
			"port": 2581,
			"host": "127.0.0.1",
		})

	if err := engine.LoadInEnd(grpcInend); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	//
	dev1 := typex.NewDevice(typex.GENERIC_PROTOCOL,
		"UART", "UART", "UART", map[string]interface{}{
			"commonConfig": map[string]interface{}{
				"frequency":   5,
				"autoRequest": true,
				"transport":   "rs485rawserial",
				"waitTime":    10,
				"timeout":     10,
			},
			"uartConfig": map[string]interface{}{
				"baudRate": 9600,
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"uart":     "COM5",
			},
			"deviceConfig": map[string]interface{}{
				"get_uuid": map[string]interface{}{
					"name":        "get_uuid",
					"description": "获取UUID",
					"protocol": map[string]interface{}{
						"in":  "AABBCCDDEEFF",
						"out": "112233445566",
					},
				},
			},
		})
	dev1.UUID = "dev1"
	if err := engine.LoadDevice(dev1); err != nil {
		t.Fatal("dev1 load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"test",
		"test",
		[]string{grpcInend.UUID},
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
		function(data)
			print("received ======> : ",data)
			local n, err = rulexlib:WriteDevice("dev1", 0, "get_uuid")
			print("WriteDevice <====== : ", n, err )
			return true, data
		end
	}
`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := rulexrpc.NewRulexRpcClient(conn)

	resp, err := client.Work(context.Background(), &rulexrpc.Data{
		Value: `get_uuid`,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
	time.Sleep(20 * time.Second)
	engine.Stop()
}
