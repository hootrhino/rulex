package test

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"
	"github.com/hootrhino/rulex/plugin/http_server/model"

	"github.com/hootrhino/rulex/typex"
)

var _DataToHttp_luaCase = `function Main(arg) for i = 1, 3, 1 do local err = applib:DataToHttp('httpServer',applib:T2J({temp = 20,humi = 13.45})) applib:log('result =>') time:Sleep(100) end return 0 end`

func Test_DataToHttp(t *testing.T) {
	RmUnitTestDbFile(t)

	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal(err)
	}
	go _start_simple_http_server(t)
	//

	httpServer := typex.NewOutEnd(typex.HTTP_TARGET,
		"HTTP", "HTTP", map[string]interface{}{
			"url": "http://127.0.0.1:8899",
			"headers": map[string]interface{}{
				"secret": "test-ok",
			},
		},
	)
	httpServer.UUID = "httpServer"
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(httpServer, ctx, cancelF); err != nil {
		t.Fatal(err)
	}

	uuid := _createTestApp(t)
	time.Sleep(1 * time.Second)
	_updateTestApp(t, uuid)

	time.Sleep(20 * time.Second)
	_deleteTestApp(t, uuid)
	engine.Stop()
}

//--------------------------------------------------------------------------------------------------
// 起一个HTTP服务器
//--------------------------------------------------------------------------------------------------

func index(w http.ResponseWriter, req *http.Request) {
	body := [1000]byte{}
	req.Body.Read(body[:req.ContentLength])
	log.Println("[OUT]        Body ===> ", string(body[:req.ContentLength]))
	log.Println("[OUT]        Header ==>", req.Header["Secret"])
	w.Write([]byte("OK"))
}

func _start_simple_http_server(t *testing.T) {
	t.Log("_start_simple_http_server")
	http.HandleFunc("/", index)
	http.ListenAndServe(":8899", nil)
}

//--------------------------------------------------------------------------------------------------
// 资源操作
//--------------------------------------------------------------------------------------------------

func _createTestApp(t *testing.T) string {
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
func _updateTestApp(t *testing.T, uuid string) {
	body := `{"uuid": "%s","name": "testlua11","version": "1.0.1","autoStart": true,"luaSource":"AppNAME='OK1'\nAppVERSION='0.0.3'\n%s"}`

	t.Logf(body, uuid, _DataToHttp_luaCase)
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", fmt.Sprintf(body, uuid, _DataToHttp_luaCase)).Output()
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
	_startTestApp(t, mApp[0].UUID)
}
func _deleteTestApp(t *testing.T, uuid string) {
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
func _startTestApp(t *testing.T, uuid string) {
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app/start?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_startTestApp: ", string(output))
}
