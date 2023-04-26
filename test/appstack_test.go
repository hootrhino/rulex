package test

import (
	"testing"
	"time"

	"github.com/hootrhino/rulex/appstack"
	"github.com/hootrhino/rulex/typex"
)

// go test -timeout 30s -run ^Test_appStack github.com/hootrhino/rulex/test -v -count=1
func Test_appStack(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()
	as := appstack.NewAppStack(engine)
	err := as.LoadApp(typex.NewApplication("test-uuid-1", "test-name",
		"1.0.1", "./apps/hello_world.lua"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(as)
	time.Sleep(10 * time.Second)
	engine.Stop()
}
