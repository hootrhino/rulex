package test

import (
	"rulex/core"
	"rulex/engine"

	"rulex/typex"
	"testing"
	"time"
)

func Test_ABS_device1(t *testing.T) {
	engine := engine.NewRuleEngine(core.InitGlobalConfig("conf/rulex.ini"))
	engine.Start()
	demoDev := &typex.Device{
		UUID:         "Test1",
		Name:         "Test1",
		Type:         "SIMPLE",
		ActionScript: "******",
		Description:  "Test1",
		Config: map[string]interface{}{
			"K": "V",
		},
	}

	engine.LoadDevice(demoDev)
	t.Log(engine.SnapshotDump())

	time.Sleep(20 * time.Second)
	engine.Stop()
}
func Test_ABS_device2(t *testing.T) {
	engine := engine.NewRuleEngine(core.InitGlobalConfig("conf/rulex.ini"))
	engine.Start()
	demoDev := &typex.Device{
		UUID:         "Test1",
		Name:         "Test1",
		Type:         "NO-SUCH",
		ActionScript: "******",
		Description:  "Test1",
		Config: map[string]interface{}{
			"K": "V",
		},
	}

	if err := engine.LoadDevice(demoDev); err != nil {
		t.Log(err)
	}
	time.Sleep(1 * time.Second)
	engine.Stop()
}
