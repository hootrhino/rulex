package test

import (
	"encoding/json"
	"os"

	"testing"

	"github.com/hootrhino/rulex/common"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

//	{
//	    "host":"127.0.0.1",
//	    "port":1800,
//	    "rack":0,
//	    "slot":1,
//	    "model":"S1200",
//	    "timeout":5,
//	    "idleTimeout":5,
//	    "readFrequency":5,
//	    "blocks":[
//	        {
//	            "tag":"V1",
//	            "address":1,
//	            "start":1,
//	            "size":10
//	        }
//	    ]
//	}
func Test_gen_config(t *testing.T) {

}
func Test_parse_config(t *testing.T) {
	config := map[string]interface{}{
		"host":          "127.0.0.1",
		"port":          1800,
		"rack":          0,
		"slot":          1,
		"model":         "S1200",
		"timeout":       5,
		"idleTimeout":   5,
		"readFrequency": 5,
		"blocks": []map[string]interface{}{
			{
				"tag":     "V1",
				"address": 1,
				"start":   1,
				"size":    10,
			},
		},
	}
	configMain := common.S1200Config{}
	configBytes, err0 := json.Marshal(&config)
	if err0 != nil {
		t.Fatal(err0)
	}
	t.Log(string(configBytes))
	if err1 := json.Unmarshal(configBytes, &configMain); err1 != nil {
		t.Fatal(err1)
	}

}

/*
*
* 测试RULEX加载 SIEMENS_PLC
*
 */
func Test_RULEX_WITH_SIEMENS_PLC(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	SIEMENS_PLC := typex.NewDevice(typex.SIEMENS_PLC,
		"PLC工站系统", "PLC工站系统", map[string]interface{}{
			"host":          "127.0.0.1",
			"port":          1800,
			"rack":          0,
			"slot":          1,
			"model":         "S1200",
			"timeout":       5,
			"idleTimeout":   5,
			"readFrequency": 5,
			"blocks": []map[string]interface{}{
				{
					"tag":     "V1",
					"address": 1,
					"start":   1,
					"size":    10,
				},
				{
					"tag":     "V2",
					"address": 1,
					"start":   1,
					"size":    10,
				},
			},
		},
	)
	SIEMENS_PLC.UUID = "SIEMENS_PLC"
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadDeviceWithCtx(SIEMENS_PLC, ctx, cancelF); err != nil {
		t.Error("SIEMENS_PLC load failed:", err)
	}
	//
	// 透传到内部EMQX
	//
	EMQX_BROKER := typex.NewOutEnd(typex.MQTT_TARGET,
		"内网MQTT桥接",
		"内网MQTT桥接", map[string]interface{}{
			"Host":      "emqx.dev.com",
			"Port":      1883,
			"DataTopic": "iothub/upstream/YK0801",
			"ClientId":  "YK0801",
			"Username":  "YK0801",
			"Password":  "YK0801",
		},
	)
	EMQX_BROKER.UUID = "EMQX_BROKER"
	ctx1, cancelF1 := typex.NewCCTX()
	if err := engine.LoadOutEndWithCtx(EMQX_BROKER, ctx1, cancelF1); err != nil {
		t.Error("mqttOutEnd load failed:", err)
	}
	// 	// 加载一个规则
	rule1 := typex.NewRule(engine,
		"uuid",
		"FROM TENCENT_IOT_INEND",
		"FROM TENCENT_IOT_INEND",
		[]string{EMQX_BROKER.UUID}[0],
		"",
		`function Success()end`,
		`
	Actions = {
		function(args)
		    rulexlib:log('EMQX_BROKER: ', data)
			return true, args
		end
	}`, `function Failed(error) print("[EMQX_BROKER Failed Callback]", error) end`)
	if err := engine.LoadRule(rule1); err != nil {
		t.Error(err)
	}

	glogger.GLogger.Warn("Received stop signal:", s)
	engine.Stop()
	os.Exit(0)
}
