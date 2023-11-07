package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"github.com/adrianmo/go-nmea"
	"github.com/hootrhino/rulex/typex"
)

// go test -timeout 30s -run ^Test_AIS_SEND_PACKET github.com/hootrhino/rulex/test -v -count=1

func Test_AIS_SEND_PACKET(t *testing.T) {
	s1 := `\1G1:370208949,g:1,s:ABC,t:2320,c:1660780800,d:110*72\!ABVDM,1,1,5,B,H69EvShlTpID@TpMUG3COOL0000,2*14`
	s, err := nmea.Parse(s1)
	if err != nil {
		log.Fatal(err)
	}
	if s.DataType() == nmea.TypeVDM {
		fmt.Println(s.String())
		m := s.(nmea.VDMVDO)
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}

/*
*
* AIS 数据发送模拟器
*
 */
func ais_sender_emulator_udp() {
	// Server address
	serverAddr := "localhost:6005"

	// Create UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Failed to resolve server address:", err)
		return
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	// Send data to the server
	s1 := `\1G1:370208949,g:1,s:ABC,t:2320,c:1660780800,d:110*72\!ABVDM,1,1,5,B,H69EvShlTpID@TpMUG3COOL0000,2*14`
	_, err = conn.Write([]byte(s1))
	if err != nil {
		fmt.Println("Failed to send data to server:", err)
		return
	}

	// Receive response from the server
	buffer := make([]byte, 1024)
	bytesRead, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Failed to receive response from server:", err)
		return
	}

	response := string(buffer[:bytesRead])
	fmt.Println("Response from server:", response)
}

func ais_sender_emulator_tcp() {
	// Connect to the server
	s1 := `!ABVDM,1,1,5,B,H69EvShlTpID@TpMUG3COOL0000,2*14`
	conn, err := net.Dial("tcp", "localhost:6005")
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(s1))
	if err != nil {
		fmt.Println("Failed to send data:", err)
		return
	}

	// Receive response from the server
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Failed to receive response:", err)
		return
	}

	response := string(buffer[:bytesRead])
	fmt.Println("Response from server:", response)
}

// go test -timeout 30s -run ^Test_generic_ais_txrx_device github.com/hootrhino/rulex/test -v -count=1
func Test_generic_ais_txrx_device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal(err)
	}
	GENERIC_AIS_RECEIVER := typex.NewDevice(typex.GENERIC_AIS_RECEIVER,
		"GENERIC_AIS_RECEIVER", "GENERIC_AIS_RECEIVER", map[string]interface{}{
			"host": "0.0.0.0",
			"port": 6005,
		})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(GENERIC_AIS_RECEIVER, ctx, cancelF); err != nil {
		t.Fatal(err)
	}
	time.Sleep(25 * time.Second)
	engine.Stop()
}
