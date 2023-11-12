# TSS200
## 简介

TC-S200 系列空气质量监测仪内置 PM2.5、TVOC、甲醛、CO2，温湿度等高精度传感器套件，可通过吸顶式或壁挂安装，RS-485 接口通过 Modbus-RTU 协议进行数据输出，通过网关组网，或配合联动模块可以用于新风联动控制。

## 配置

```json
{
    "name": "RTU485_THER",
    "type": "RTU485_THER",
    "config": {
        "mode": "UART",
        "timeout": 10,
        "frequency": 5,
        "autoRequest": true,
        "config": {
            "uart": "/dev/ttyUSB0",
            "dataBits": 8,
            "parity": "N",
            "stopBits": 1,
            "baudRate": 9600,
            "ip": "127.0.0.1",
            "port": 502
        },
        "registers": [
            {
                "tag": "node1",
                "function": 3,
                "slaverId": 1,
                "address": 0,
                "quantity": 2
            }
        ]
    },
    "description": "RTU485_THER"
}
```

## 数据

### 读
#### 1. CMD
无指令。
#### 1. Args
无参数。

#### 3. 数据样例
```json
{
	float32 `json:"temp"` //系数: 0.01
	float32 `json:"hum"`  //系数: 0.01
	uint16  `json:"pm1"`
	uint16  `json:"pm25"`
	uint16  `json:"pm10"`
	uint16  `json:"co2"`
	float32 `json:"tvoc"` //系数: 0.001
	float32 `json:"choh"` //系数: 0.001
	float32 `json:"eco2"` //系数: 0.001
}
```

### 写
#### 1. CMD
无指令。

#### 2. Args
无参数。

#### 3. 数据样例
无写数据。

## 案例

```lua
-- Success
function Success()
    print("Success")
end
-- Failed
function Failed(error)
    print("Error:", error)
end

-- Actions
Actions = {
    function(args)
        local _, err = rulexlib:ReadDevice(device, 0, "all")
        if (err ~= nil) then
            return false, data
        end
        return true, args
    end
}

```

## 说明

-  暂无

## 社区

### 维护者

- wwhai
