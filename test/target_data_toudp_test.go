package test

import (
	"fmt"
	"log"
	"os/exec"

	"net"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/plugin/http_server/model"

	"github.com/hootrhino/rulex/typex"
)

var _DataToUdp_luaCase = `function Main(arg) for i = 1, 3, 1 do local err = applib:DataToUdp('UdpServer',applib:T2J({temp = 20,humi = 13.45})) applib:log('result =>',err) time:Sleep(100) end return 0 end`

// go test -timeout 30s -run ^Test_DataToUdp github.com/hootrhino/rulex/test -v -count=1

func Test_DataToUdp(t *testing.T) {
	RmUnitTestDbFile(t)

	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// UdpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal(err)
	}
	go _start_simple_Udp_server()
	//

	UdpServer := typex.NewOutEnd(typex.UDP_TARGET,
		"Udp", "Udp", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 8891,
		},
	)
	UdpServer.UUID = "UdpServer"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF

	if err := engine.LoadOutEndWithCtx(UdpServer, ctx1, cancelF1); err != nil {
		t.Fatal(err)
	}

	uuid := _createTestApp_1(t)
	time.Sleep(1 * time.Second)
	_updateTestApp_1(t, uuid)

	time.Sleep(20 * time.Second)
	_deleteTestApp_1(t, uuid)
	engine.Stop()
}

//--------------------------------------------------------------------------------------------------
// 起一个Udp服务器
//--------------------------------------------------------------------------------------------------

func _start_simple_Udp_server() {

	//UDP服务器监听 接收广播数据
	udp_addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8891")
	if err != nil {
		log.Fatal(err)
	}
	udp_conn, err := net.ListenUDP("udp", udp_addr)
	if err != nil {
		log.Fatal(err)
	}
	defer udp_conn.Close()

	data := [100]byte{}
	for {
		udp_conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := udp_conn.ReadFromUDP(data[:])
		if err != nil {
			log.Fatal(err)
		}
		udp_conn.SetReadDeadline(time.Time{})

		log.Println("UDP Received ============>:", string(data[:n]))
	}

}

//--------------------------------------------------------------------------------------------------
// 资源操作
//--------------------------------------------------------------------------------------------------

func _createTestApp_1(t *testing.T) string {
	// 通过接口创建一个App
	body := `{"name": "testlua1","version": "1.0.0","autoStart": false,"description": "hello world"}`
	output, err := exec.Command("curl",
		"-X", "POST", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", body).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_createApp: ", string(output))
	//
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 1, len(mApp))
	t.Log(mApp[0].UUID)
	assert.Equal(t, mApp[0].Name, "testlua1")
	assert.Equal(t, mApp[0].Version, "1.0.0")
	assert.Equal(t, mApp[0].AutoStart, false)
	return mApp[0].UUID
}
func _updateTestApp_1(t *testing.T, uuid string) {
	body := `{"uuid": "%s","name": "testlua11","version": "1.0.1","autoStart": true,"luaSource":"AppNAME='OK1'\nAppVERSION='0.0.3'\n%s"}`

	t.Logf(body, uuid, _DataToUdp_luaCase)
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", fmt.Sprintf(body, uuid, _DataToUdp_luaCase)).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_updateApp: ", string(output))
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 1, len(mApp))
	t.Log("APP UUID ==> ", mApp[0].UUID)
	assert.Equal(t, mApp[0].Name, "testlua11")
	assert.Equal(t, mApp[0].Version, "1.0.1")
	assert.Equal(t, mApp[0].AutoStart, true)
	// _startTestApp
	time.Sleep(1 * time.Second)
	_startTestApp_1(t, mApp[0].UUID)
}
func _deleteTestApp_1(t *testing.T, uuid string) {
	// 删除一个App
	output, err := exec.Command("curl", "-X", "DELETE", "http://127.0.0.1:2580/api/v1/app?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_deleteApp: ", string(output))
	//
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 0, len(mApp))
}
func _startTestApp_1(t *testing.T, uuid string) {
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app/start?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_startTestApp: ", string(output))
}
