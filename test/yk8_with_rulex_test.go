package test

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

//
// 下行请求 ----------------------------------------------------------------------------------------
// -- 应用调用设备行为 Topic: $thing/down/action/{ProductID}/{DeviceName}
// -- 设备响应行为执行结果 Topic: $thing/up/action/{ProductID}/{DeviceName}
// -------------------------------------------------------------------------------------------------
// {
//     "method":"action",
//     "clientToken":"111111111111111111",
//     "actionId":"control",
//     "timestamp":1657105901,
//     "params":{
//         "sw1":0,
//         "sw2":0,
//         "sw3":0,
//         "sw4":0,
//         "sw5":0,
//         "sw6":0,
//         "sw7":0,
//         "sw8":0
//     }
// }
// 下行请求设备的回复 ------------------------------------------------------------------------------
// {
//     "method": "action_reply",
//     "clientToken": "111111111111111111",
//     "code": 0,
//     "status": "OK",
//     "response": {
//         // 定义的结构体
//      }
// }
// 下行控制 ----------------------------------------------------------------------------------------
// -- 下发 Topic: $thing/down/property/{ProductID}/{DeviceName}
// -- 响应 Topic: $thing/up/property/{ProductID}/{DeviceName}
// -------------------------------------------------------------------------------------------------
// {
//     "method":"control",
//     "clientToken":"clientToken-447ec5b5-a978-477f-aaa7-99a1c78e342a",
//     "params":{
//         "switchers":{
//             "sw7":0,
//             "sw8":0,
//             "sw1":0,
//             "sw2":0,
//             "sw3":0,
//             "sw4":0,
//             "sw5":0,
//             "sw6":0
//         }
//     }
// }
// 事件上报 ----------------------------------------------------------------------------------------
// -- 设备事件上行请求 Topic： $thing/up/event/{ProductID}/{DeviceName}
// -- 设备事件上行响应 Topic： $thing/down/event/{ProductID}/{DeviceName}
// -------------------------------------------------------------------------------------------------
// {
//    "method":"event_post",
//    "clientToken":"123",
//    "version":"1.0",
//    "eventId":"PowerAlarm",
//    "type":"fault",
//    "timestamp":11111111,
//    "params":{
//        "Voltage":2.8,
//        "Percent":20
//    }
// }
//
/*
*
* 测试RULEX加载Yk8
*
 */
func Test_RULEX_WITH_YK08(t *testing.T) {
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	glogger.StartGLogger(core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	core.SetLogLevel()
	core.SetPerformance()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	engine := engine.NewRuleEngine(mainConfig)
	engine.Start()

	//
	// RTU485_THER Inend
	//
	// YK801 := typex.NewDevice(typex.YK08_RELAY,
	// 	"继电器控制器", "继电器控制器", "", map[string]interface{}{
	// 		"slaverIds": []uint8{1},
	// 		"timeout":   5,
	// 		"frequency": 5,
	// 		"config": map[string]interface{}{
	// 			"uart":     "COM3",
	// 			"baudRate": 9600,
	// 			"dataBits": 8,
	// 			"parity":   "N",
	// 			"stopBits": 1,
	// 		},
	// 	})
	// YK801.UUID = "YK801"
	// if err := engine.LoadDevice(YK801); err != nil {
	// 	t.Error("YK08_RELAY load failed:", err)
	// }
	//
	// 腾讯云MQTT
	//
	TENCENT_IOT_INEND := typex.NewInEnd(
		typex.TENCENT_IOT_HUB,
		"TENCENT_IOT_HUB",
		"TENCENT_IOT_HUB", map[string]interface{}{
			"Host":       "Y0ST19XLP1.iotcloud.tencentdevices.com",
			"Port":       1883,
			"productId":  "Y0ST19XLP1",
			"deviceName": "YK8_001",
			"ClientId":   "Y0ST19XLP1YK8_001",
			"Username":   "Y0ST19XLP1YK8_001;12010126;Y0YVU;1657592838",
			"Password":   "b679f531ca4eacbd280c87a6d027cd6aba7d63c0e2f1310fd4ec6e31d2fe7163;hmacsha256",
		},
	)
	TENCENT_IOT_INEND.UUID = "TENCENT_IOT_INEND"
	if err := engine.LoadInEnd(TENCENT_IOT_INEND); err != nil {
		t.Error("TENCENT_IOT_INEND load failed:", err)
	}
	//
	// 透传到内部平台
	//
	// 	mqttOutEnd := typex.NewOutEnd(typex.MQTT_TARGET,
	// 		"内网MQTT桥接",
	// 		"内网MQTT桥接", map[string]interface{}{
	// 			"Host":      "emqx.dev.inrobot.cloud",
	// 			"Port":      1883,
	// 			"DataTopic": "iothub/upstream/YK0801",
	// 			"ClientId":  "YK0801",
	// 			"Username":  "YK0801",
	// 			"Password":  "YK0801",
	// 		},
	// 	)
	// 	mqttOutEnd.UUID = "mqttOutEnd"
	// 	if err := engine.LoadOutEnd(mqttOutEnd); err != nil {
	// 		t.Error("mqttOutEnd load failed:", err)
	// 	}
	// 	// 加载一个规则
	rule1 := typex.NewRule(engine,
		"uuid",
		"FROM TENCENT_IOT_INEND",
		"FROM TENCENT_IOT_INEND",
		[]string{TENCENT_IOT_INEND.UUID},
		[]string{},
		`function Success()end`,
		`
	Actions = {
		function(data)
		    rulexlib:log('TENCENT_IOT_INEND: ', data)
			return true, data
		end
	}`, `function Failed(error) print("[TENCENT_IOT_INEND Failed Callback]", error) end`)
	// 	rule2 := typex.NewRule(engine,
	// 		"uuid",
	// 		"数据推送至IOTHUB",
	// 		"数据推送至IOTHUB",
	// 		[]string{},
	// 		[]string{YK801.UUID},
	// 		`function Success()end`,
	// 		`
	// Actions = {
	// 	function(data)
	// 		rulexlib:log(data)
	// 		return true, data
	// 	end
	// }`, `function Failed(error) print("[YK801 Failed Callback]", error) end`)
	if err := engine.LoadRule(rule1); err != nil {
		t.Error(err)
	}
	// if err := engine.LoadRule(rule2); err != nil {
	// 	t.Error(err)
	// }
	s := <-c
	glogger.GLogger.Warn("Received stop signal:", s)
	engine.Stop()
	os.Exit(0)
}
