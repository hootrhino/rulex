package test

import (
	"encoding/json"
	"fmt"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"log"
	"testing"
	"time"

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
func ais_sender_emulator() {

}

// go test -timeout 30s -run ^Test_generic_ais_txrx_device github.com/hootrhino/rulex/test -v -count=1
func Test_generic_ais_txrx_device(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal(err)
	}
	GENERIC_AIS := typex.NewDevice(typex.GENERIC_AIS,
		"GENERIC_AIS", "GENERIC_AIS", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 9980,
		})

	if err := engine.LoadDevice(GENERIC_AIS); err != nil {
		t.Fatal(err)
	}
	time.Sleep(25 * time.Second)
	engine.Stop()
}
