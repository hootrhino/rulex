package test

import (
	"testing"
	"time"

	"github.com/i4de/rulex/appstack"
)

// go test -timeout 30s -run ^Test_appStack github.com/i4de/rulex/test -v -count=1
func Test_appStack(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()
	as := appstack.NewAppStack(engine)
	err := as.LoadApp(appstack.NewApplication("test-uuid-1", "test-name",
		"1.0.1", "helloworld_1.0.0.lua"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(as)
	time.Sleep(30 * time.Second)
	engine.Stop()
}
