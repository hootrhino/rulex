## 简介
USR-G776 是有人第二代 4G DTU 产品，支持移动，联通，电信 4G 和移动，联通3G和2G 网络制式，以“透
传”作为功能核心，高度易用性，用户可方便快速的集成于自己的系统中。该 DTU 软件功能完善，覆盖绝大多
数常规应用场景，用户只需通过简单的设置，即可实现串口到网络的双向数据透明传输。

RULEX对USR-G776做了一个非常简单的支持：直接向串口透传数据，因此在使用该设备之前最好是配置好。

官方文档：https://www.usr.cn/Download/806.html
## 配置
```go
type _G776CommonConfig struct {
	Tag         string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	Frequency   int64  `json:"frequency" validate:"required" title:"采集频率"`
	AutoRequest bool   `json:"autoRequest" title:"启动轮询"`
}
type CommonUartConfig struct {
	Timeout  int    `json:"timeout" validate:"required"`
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
}
type _G776Config struct {
	CommonConfig _G776CommonConfig       `json:"commonConfig" validate:"required"`
	UartConfig   common.CommonUartConfig `json:"uartConfig" validate:"required"`
}

```

## 数据处理示例

```lua
-- Success
function Success()

end

-- Failed
function Failed(error)
    print("failed:", error)
end

-- Actions

Actions = {
    function(args)
        local dataTable, err1 = rulexlib:J2T(data)
        if err1 ~= nil then
            return true, args
        end
        for _k, entity in pairs(dataTable) do
            data:ToUsrG776DTU("uuid", "DATA" ,rulexlib:T2J(entity["value"]))
        end
        return true, args
    end
}

```
### 注意事项
第二个参数必须是“DATA”，第三个参数可以是字符串、十六进制等.