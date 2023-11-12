## SNMP协议
SNMP（Simple Network Management Protocol）是一种用于网络管理的协议。它被广泛用于监视和管理网络中的设备、系统和应用程序。SNMP允许网络管理员通过网络收集和组织设备的管理信息，以便实时监控网络的状态、性能和健康状况。

SNMP基于客户端-服务器模型，其中有两个主要组件：

1. 管理器（Manager）：管理器是网络管理员或管理系统的一部分，用于监视和控制网络设备。它通过SNMP协议发送请求并接收响应，从被管理设备中获取信息并采取必要的操作。

2. 代理（Agent）：代理是安装在被管理设备上的软件模块，负责收集和维护设备的管理信息，并响应来自管理器的请求。代理将设备的状态、配置和性能信息以适当的格式暴露给管理器。

SNMP使用一组标准的管理信息库（Management Information Base，MIB）来定义设备和系统的管理信息。MIB是一个层次化的数据库，包含了设备的各种参数、统计信息和配置设置。管理器可以通过SNMP协议向代理发送请求，例如获取特定参数的值、设置参数的值、触发操作等。

SNMP协议支持各种版本，其中最常用的是SNMPv1、SNMPv2c和SNMPv3。每个版本都具有不同的功能和安全性特性，以适应不同的网络管理需求。
## 设备配置
```go

type _SNMPCommonConfig struct {
	AutoRequest bool  `json:"autoRequest" title:"启动轮询"`
	Frequency   int64 `json:"frequency" validate:"required" title:"采集频率"`
}
type GenericSnmpConfig struct {
	// Target is an ipv4 address.
	Target string `json:"target" validate:"required" title:"Target" info:"Target"`
	// Port is a port.
	Port uint16 `json:"port" validate:"required" title:"Port" info:"Port"`
	// Transport is the transport protocol to use ("udp" or "tcp"); if unset "udp" will be used.
	Transport string `json:"transport" validate:"required" title:"Transport" info:"Transport"`
	// Community is an SNMP Community string.
	Community string `json:"community" validate:"required" title:"Community" info:"Community"`
}

type _GSNMPConfig struct {
	CommonConfig _SNMPCommonConfig        `json:"commonConfig" validate:"required"`
	SNMPConfig   common.GenericSnmpConfig `json:"snmpConfig" validate:"required"`
}

```
## 数据示例
```json
{
    "PCHost":"127.0.0.1",
    "PCDescription":"Linux x86_64",
    "PCUserName":"demo",
    "PCHardIFaces":[],
    "PCTotalMemory":0
}
```
## 数据解析示例
```lua

function (data)
    local DataT, err = rulexlib:J2T(data)
    if err ~= nil then
        return true, args
    end
    -- Do your business
    rulexlib:log(DataT['PCHost'])
    rulexlib:log(DataT['PCDescription'])
    rulexlib:log(DataT['PCUserName'])
    rulexlib:log(DataT['PCHardIFaces'])
    rulexlib:log(DataT['PCTotalMemory'])
end

```