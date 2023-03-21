package test

import (
	"time"

	httpserver "github.com/i4de/rulex/plugin/http_server"

	"testing"

	"github.com/i4de/rulex/typex"
)

/*
*
* Test 485 sensor gateway
*
 */

func Test_modbus_485_yk8(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", httpserver.NewHttpApiServer()); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}

	// YK8 Inend
	YK8Device := typex.NewDevice(typex.YK08_RELAY,
		"继电器", "继电器", "继电器", map[string]interface{}{
			"timeout":   5,
			"frequency": 5,
			"config": map[string]interface{}{
				"uart":     "COM4",
				"dataBits": 8,
				"parity":   "N",
				"stopBits": 1,
				"baudRate": 9600,
			},
			"registers": []map[string]interface{}{
				{
					"tag":      "yk08-001",
					"function": 3,
					"slaverId": 3,
					"address":  0,
					"quantity": 1,
				},
			},
		})
	YK8Device.UUID = "YK8Device1"
	if err := engine.LoadDevice(YK8Device); err != nil {
		t.Fatal("YK8Device load failed:", err)
	}

	tencentIothub := typex.NewInEnd(typex.TENCENT_IOT_HUB,
		"MQTT", "MQTT", map[string]interface{}{
			"host":       "10.55.16.144",
			"port":       1883,
			"clientId":   "d021c10445c3142af979e4decdd22a797",
			"username":   "d021c10445c3142af979e4decdd22a797",
			"password":   "da2b04e6b03c4aa587b3c43486e6ab2d",
			"productId":  "pe6a5a5889f2b449d97110965fa95e91f",
			"deviceName": "d021c10445c3142af979e4decdd22a797",
		})
	tencentIothub.UUID = "tencentIothub"

	if err := engine.LoadInEnd(tencentIothub); err != nil {
		t.Fatal("mqttOutEnd load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"数据推送至IOTHUB",
		"数据推送至IOTHUB",
		[]string{tencentIothub.UUID}, // 数据来自MQTT Server
		[]string{},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
Actions = {
	function (data)
	print('Received data from iothub:', data)
	local source = 'tencentIothub'
	local device = 'YK8Device1'
	local dataT, err = rulexlib:J2T(data)
	if (err ~= nil) then
		print('Received data from iothub parse to json error:', err)
		return false, data
	end
	-- Action
	if dataT['method'] == 'action' then
		local actionId = dataT['actionId']
		if actionId == 'get_state' then
			local readData, err = rulexlib:ReadDevice(0, device)
			if (err ~= nil) then
				print('ReadDevice data from device error:', err)
				return false, data
			end
			print('ReadDevice data from device:', readData)
			local readDataT, err = rulexlib:J2T(readData)
			if (err ~= nil) then
				print('Parse ReadDevice data from device error:', err)
				return false, data
			end
			local yk08001State = readDataT['yk08-001']
			print('yk08001State:', yk08001State['value'])
			local _, err = iothub:ActionSuccess(source, dataT['id'], yk08001State['value'])
			if (err ~= nil) then
				print('ActionReply error:', err)
				return false, data
			end
		end
	end
	-- Property
	if dataT['method'] == 'property' then
		local schemaParams = rulexlib:J2T(dataT['data'])
		print('schemaParams:', schemaParams)
		local n1, err = rulexlib:WriteDevice(device, 0, rulexlib:T2J({{
			['function'] = 15,
			['slaverId'] = 3,
			['address'] = 0,
			['quantity'] = 1,
			['value'] = rulexlib:T2Str({
				[1] = schemaParams['sw8'],
				[2] = schemaParams['sw7'],
				[3] = schemaParams['sw6'],
				[4] = schemaParams['sw5'],
				[5] = schemaParams['sw4'],
				[6] = schemaParams['sw3'],
				[7] = schemaParams['sw2'],
				[8] = schemaParams['sw1']
			})
		}}))
		if (err ~= nil) then
			print('WriteDevice error:', err)
			local _, err = iothub:PropertyFailed(source, dataT['id'])
			if (err ~= nil) then
				print('Reply error:', err)
				return false, data
			end
		else
			local _, err = iothub:PropertySuccess(source, dataT['id'], {})
			if (err ~= nil) then
				print('Reply error:', err)
				return false, data
			end
		end
	end
	return true, data
end
}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Second)
	engine.Stop()
}
