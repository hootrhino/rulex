## Http 采集器
主要用来请求远程HTTP接口，采集这类设备的数据。

## 配置

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "uuid": "DEVICEP6VNYAAK",
        "gid": "DROOT",
        "name": "HTTP请求设备",
        "type": "GENERIC_HTTP_DEVICE",
        "state": 1,
        "config": {
            "commonConfig": {
                "autoRequest": true,
                "frequency": 1000,
                "timeout": 3000
            },
            "httpConfig": {
                "headers": {
                    "token": "12345"
                },
                "url": "http://127.0.0.1:8080"
            }
        },
        "description": ""
    }
}
```

## 测试
下面是一个传感器采集到的数据:
```lua
{
    "device_id": 1, // 设备的唯一标识符
    "recv_time": "2023-12-01T15:27:25+08:00", // 数据接收的时间戳
    "bat_voltage": 0, // 电池电压，单位通常为伏特
    "longitude": 0, // 经度，表示设备所在位置的东经或西经度数
    "latitude": 0, // 纬度，表示设备所在位置的北纬或南纬度数
    "air_height": 0, // 距离海平面的空气高度，单位通常为米
    "water_temp": 18.91, // 水温，单位通常为摄氏度
    "salinity": 627.317, // 盐度，表示水中溶解盐的含量
    "dissolved_oxygen": 0, // 溶解氧含量，单位通常为毫克/升
    "ph_value": 5.95799, // pH值，表示水质的酸碱程度
    "wind_speed": 4.48, // 风速，单位通常为米/秒
    "wind_direction": 37, // 风向，通常以北为0度，顺时针增加
    "air_temp": 13.9, // 空气温度，单位通常为摄氏度
    "air_pressure": 102.6, // 空气压力，单位通常为百帕
    "air_humidity": 63.9, // 空气湿度，通常以百分比表示
    "noise": 42, // 噪音水平，单位通常为分贝
    "wave_height": 0, // 波高，单位通常为米
    "mean_wave_period": 0, // 平均波周期，单位通常为秒
    "peak_wave_period": 0, // 最大波周期，单位通常为秒
    "mean_wave_direction": 0 // 平均波向，通常以北为0度，顺时针增加
}

```

使用LUA脚本解析并且打印出来:

```lua
Actions = {
	function(args)
		local jsonTable = {
			device_id = "设备唯一标识符",
			recv_time = "数据接收时间戳",
			bat_voltage = "电池电压",
			longitude = "经度",
			latitude = "纬度",
			air_height = "空气高度",
			water_temp = "水温",
			salinity = "盐度",
			dissolved_oxygen = "溶解氧含量",
			ph_value = "pH值",
			wind_speed = "风速",
			wind_direction = "风向",
			air_temp = "空气温度",
			air_pressure = "空气压力",
			air_humidity = "空气湿度",
			noise = "噪音水平",
			wave_height = "波高",
			mean_wave_period = "平均波周期",
			peak_wave_period = "最大波周期",
			mean_wave_direction = "平均波向"
		}
		local dataT = json:J2T(args)
		for k, v in pairs(dataT) do
			stdlib:Debug(jsonTable[k] .. ": " .. v)
		end
		return true, args
	end
}
```