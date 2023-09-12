package test

import (
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/engine"

	"testing"
	"time"

	"github.com/hootrhino/rulex/typex"
)

func Test_ABS_device1(t *testing.T) {
	engine := engine.InitRuleEngine(core.InitGlobalConfig("conf/rulex.ini"))
	engine.Start()
	demoDev := &typex.Device{
		UUID:        "Test1",
		Name:        "Test1",
		Type:        "SIMPLE",
		Description: "Test1",
		Config: map[string]interface{}{
			"K": "V",
		},
	}

	ctx, cancelF := typex.NewCCTX()
	engine.LoadDeviceWithCtx(demoDev, ctx, cancelF)
	t.Log(engine.SnapshotDump())

	time.Sleep(20 * time.Second)
	engine.Stop()
}
func Test_ABS_device2(t *testing.T) {
	engine := engine.InitRuleEngine(core.InitGlobalConfig("conf/rulex.ini"))
	engine.Start()
	demoDev := &typex.Device{
		UUID:        "Test1",
		Name:        "Test1",
		Type:        "NO-SUCH",
		Description: "Test1",
		Config: map[string]interface{}{
			"K": "V",
		},
	}
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(demoDev, ctx, cancelF); err != nil {
		t.Log(err)
	}
	time.Sleep(1 * time.Second)
	engine.Stop()
}
