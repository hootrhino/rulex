package test

import (
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/component/appstack"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"testing"

	"github.com/hootrhino/rulex/typex"
)

// 读出来的字节缓冲默认大小
const __DEFAULT_BUFFER_SIZE = 100

// 传输形式：
// `rawtcp`, `rawudp`, `rawserial`
const rawtcp string = "TCP"
const rawudp string = "rawudp"
const rawserial string = "UART"

type _CPDCommonConfig struct {
	Transport string `json:"transport" validate:"required"` // 传输协议
	RetryTime int    `json:"retryTime" validate:"required"` // 几次以后重启,0 表示不重启
}

/*
*
* 自定义协议
*
 */
type _CustomProtocolConfig struct {
	CommonConfig _CPDCommonConfig        `json:"commonConfig" validate:"required"`
	UartConfig   common.CommonUartConfig `json:"uartConfig" validate:"required"`
	HostConfig   common.HostConfig       `json:"hostConfig" validate:"required"`
}

/*
*
* Test_IcmpSender_Device
*
 */
// go test -timeout 30s -run ^Test_Modbus_App_Write github.com/hootrhino/rulex/test -v -count=1
func Test_Modbus_App_Write(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer(engine)); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}
	modbusDevice := typex.NewDevice(typex.GENERIC_PROTOCOL,
		"GENERIC_PROTOCOL", "GENERIC_PROTOCOL", map[string]interface{}{
			// "mode": "TCP",
			"mode":        "UART",
			"autoRequest": false,
			"timeout":     10,
			"frequency":   5,
			"config": _CustomProtocolConfig{
				CommonConfig: _CPDCommonConfig{
					Transport: rawserial,
					RetryTime: 5,
				},
				UartConfig: common.CommonUartConfig{
					Timeout:  3000,
					Uart:     "COM12",
					BaudRate: 9600,
					DataBits: 8,
					Parity:   "N",
					StopBits: 1,
				},
				HostConfig: common.HostConfig{
					Timeout: 50,
					Host:    "127.0.0.1",
					Port:    5200,
				},
			},
		})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(modbusDevice, ctx, cancelF); err != nil {
		t.Fatal(err)
	}

	luas := `
function Main(arg)
    local error1 = modbus:F5("uuid1", 0, 1, "00")
    time:Sleep(1000)
    local error2 = modbus:F5("uuid1", 0, 1, "01")
    time:Sleep(1000)
    return 0
end
`
	appstack.LoadApp(&appstack.Application{
		UUID:      "uuid1",
		Name:      "uuid1",
		Version:   "uuid1",
		AutoStart: true,
	}, luas)
	time.Sleep(20 * time.Second)
	engine.Stop()
}
