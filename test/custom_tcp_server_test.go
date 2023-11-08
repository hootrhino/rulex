package test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/appstack"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/typex"
)

/*
*
  - 自定义协议服务,一共测试5个短报文,前面表示序号，后面表示数据，无任何意义
  - 当Lua调用接口请求的时候，会返回这些报文
    // * 0x01:01
    // * 0x02:02 03 04
    // * 0x03:0A 0B 0C 0D
    // * 0x04:11 22 33 44 55
    // * 0x05:AA BB CC DD EE FF

*
*/
func StartCustomTCPServer() {
	listener, err := net.Listen("tcp", ":3399")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	fmt.Println("listening:", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 10)
	for {
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received Request From Custom TCP:", data[:n])
		if n > 0 {
			if data[0] == 1 {
				conn.Write([]byte{0x01})
			}
			if data[0] == 2 {
				conn.Write([]byte{0x02, 0x03, 0x04})
			}
			if data[0] == 3 {
				conn.Write([]byte{0x0A, 0x0B, 0x0C, 0x0D})
			}
			if data[0] == 4 {
				conn.Write([]byte{0x11, 0x22, 0x33, 0x44, 0x55})
			}
			if data[0] == 5 {
				conn.Write([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
			}
		}
	}

}

/*
*
* 模拟请求
*
 */
func CustomTCPRequestEmu() {
	conn, err := net.Dial("tcp", "127.0.0.1:3399")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	data := make([]byte, 10)

	conn.Write([]byte{1})
	fmt.Println("Send>>>>:", 1)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	conn.Read(data)
	fmt.Println("Read>>>>:", data)
	//
	conn.Write([]byte{2})
	fmt.Println("Send>>>>:", 2)
	conn.Read(data)
	fmt.Println("Read>>>>:", data)
	//
	conn.Write([]byte{3})
	fmt.Println("Send>>>>:", 3)
	conn.Read(data)
	fmt.Println("Read>>>>:", data)
	//
	conn.Write([]byte{4})
	fmt.Println("Send>>>>:", 4)
	conn.Read(data)
	fmt.Println("Read>>>>:", data)
	//
	conn.Write([]byte{5})
	fmt.Println("Send>>>>:", 5)
	conn.Read(data)
	fmt.Println("Read>>>>:", data)
}

// // go test -timeout 30s -run ^TestCustomTCP github.com/hootrhino/rulex/test -v -count=1
// func TestCustomTCP(t *testing.T) {
// 	go StartCustomTCPServer()
// 	time.Sleep(1000 * time.Millisecond)
// 	CustomTCPRequestEmu()
// }

/*
*
* Test_data_to_http
*
 */
func TestCustomTCP(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()
	// go StartCustomTCPServer()
	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	dev1 := typex.NewDevice(typex.GENERIC_PROTOCOL,
		"UART", "UART", map[string]interface{}{
			"commonConfig": map[string]interface{}{
				"transport": "TCP",
				"retryTime": 5,
				"frequency": 100,
			},
			"hostConfig": map[string]interface{}{
				"host": "127.0.0.1",
				"port": 3399,
			},
			"uartConfig": map[string]interface{}{
				"baudRate": 9600,
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"uart":     "COM3",
				"timeout":  2000,
			},
		})
	dev1.UUID = "Test1"
	ctx1, cancel := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(dev1, ctx1, cancel); err != nil {
		t.Fatal("Test1 load failed:", err)
	}

	if err := appstack.LoadApp(appstack.NewApplication(
		"Test1", "Name", "Version"), ""); err != nil {
		t.Fatal("app Load failed:", err)
		return
	}
	if err2 := appstack.StartApp("Test1"); err2 != nil {
		t.Fatal("app Load failed:", err2)
		return
	}
	t.Log(engine.SnapshotDump())
	time.Sleep(10 * time.Second)
	engine.Stop()
}
