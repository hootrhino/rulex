package test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	serial "github.com/wwhai/tarmserial"

	"github.com/hootrhino/rulex/component/rulexrpc"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FFFFFF014CB2AA55
// go test -timeout 30s -run ^TestHexEncoding github.com/hootrhino/rulex/test -v -count=1
func TestHexEncoding(t *testing.T) {
	hexs := []byte{255, 255, 255, 1, 76, 178, 170, 85}
	s := hex.EncodeToString(hexs)
	t.Log(fmt.Sprintf("%X", hexs) == s)
	t.Log(fmt.Sprintf("%x", hexs) == s)
	t.Log(hex.DecodeString(s))
}

// go test -timeout 30s -run ^TestCheckSUM github.com/hootrhino/rulex/test -v -count=1

func TestCheckSUM(t *testing.T) {
	hexs := [8]byte{0xFF, 0xFF, 0xFF, 0x01, 0x4C, 0xB2, 0xAA, 0x55}
	for i, v := range hexs {
		t.Logf("%d %d %X", i, v, v)
	}
	checksumBegin := 0    // 校验起点
	checksumEnd := 7      // 校验结束点
	checksumValuePos := 6 // 校验比对位置
	t.Log(utils.XOR(hexs[checksumBegin:checksumEnd]) == int(hexs[checksumValuePos]))
}

// go test -timeout 30s -run ^TestCustomProtocolDevice github.com/hootrhino/rulex/test -v -count=1

func TestCustomProtocolDevice(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer(engine)); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd(typex.GRPC,
		"Rulex Grpc InEnd", "Rulex Grpc InEnd", map[string]interface{}{
			"port": 2581,
			"host": "127.0.0.1",
		})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	//
	dev1 := typex.NewDevice(typex.GENERIC_PROTOCOL,
		"UART", "UART", map[string]interface{}{
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
				"uart":     "COM3",
			},
			"deviceConfig": map[string]interface{}{
				"1": map[string]interface{}{
					"name":             "get_uuid",
					"rw":               1,
					"description":      "获取UUID",
					"bufferSize":       4,       // 期望返回几个字节
					"timeout":          1000,    // 串口读写超时
					"checksum":         "CRC16", // 校验算法
					"checksumValuePos": 7,       // 校验比对位置
					"checksumBegin":    1,       // 校验起点
					"checksumEnd":      2,       // 校验结束点
					"autoRequest":      true,    // 是否开启轮询
					"autoRequestGap":   600,     // 轮询间隔
					"protocol": map[string]interface{}{
						"in":  "FFFFFF014CB2AA55",
						"out": "FFFFFF014CB2AA55",
					},
				},
			},
		})
	dev1.UUID = "dev1"
	ctx1, cancelF1 := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(dev1, ctx1, cancelF1); err != nil {
		t.Fatal("dev1 load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"test",
		"test",
		[]string{},
		[]string{dev1.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
		function(args)
			   print("Received Inend Data ======> : ",data)
			return true, args
		end,
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

// go test -timeout 30s -run ^Test_SerialPortRW github.com/hootrhino/rulex/test -v -count=1
func Test_SerialPortRW(t *testing.T) {
	config := serial.Config{
		Name:     "COM15",
		Baud:     9600,
		Size:     8,
		Parity:   'N',
		StopBits: 1,
	}
	serialPort, err := serial.OpenPort(&config)
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := hex.DecodeString("FFFFFF014CB2AA55")
	result := [7]byte{}
	serialPort.Write((bytes))
	// time.Sleep(time.Millisecond * 60)
	n1, _ := serialPort.Read(result[:])
	t.Log("serialPort.Read 1:", n1, result[:n1])
	serialPort.Write((bytes))
	// time.Sleep(time.Millisecond * 60)
	n2, _ := serialPort.Read(result[:])
	t.Log("serialPort.Read 2:", n2, result[:n2])
	serialPort.Write((bytes))
	// time.Sleep(time.Millisecond * 60)
	n3, _ := serialPort.Read(result[:])
	t.Log("serialPort.Read 3:", n3, result[:n3])
	serialPort.Write((bytes))
	// time.Sleep(time.Millisecond * 60)
	n4, _ := serialPort.Read(result[:])
	t.Log("serialPort.Read 4:", n4, result[:n4])
	serialPort.Write((bytes))
	// time.Sleep(time.Millisecond * 60)
	n5, _ := serialPort.Read(result[:])
	t.Log("serialPort.Read 5:", n5, result[:n5])
	serialPort.Write((bytes))
	n6, _ := serialPort.Read(result[:])
	t.Log("serialPort.Read 6:", n6, result[:n6])
}
