package test

import (
	"github.com/hootrhino/rulex/glogger"
	httpserver "github.com/hootrhino/rulex/plugin/http_server"

	"testing"
	"time"

	"github.com/hootrhino/rulex/typex"
)

/*
*
* 摄像头拉流
*
 */
// go test -timeout 30s -run ^Test_Generic_camera github.com/hootrhino/rulex/test -v -count=1

func Test_Generic_camera(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
		t.Fatal(err)
	}
	GENERIC_CAMERA := typex.NewDevice(typex.GENERIC_CAMERA,
		"GENERIC_CAMERA", "GENERIC_CAMERA", map[string]interface{}{
			"maxThread":  10,
			"inputMode":  "LOCAL",
			"device":     "video0",
			"rtspUrl":    "rtsp://127.0.0.1",
			"outputMode": "JPEG_STREAM",
			"outputAddr": "0.0.0.0:2599",
		})
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(GENERIC_CAMERA, ctx, cancelF); err != nil {
		t.Fatal(err)
	}
	time.Sleep(25 * time.Second)
	engine.Stop()
}
