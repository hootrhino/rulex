# OPCUA 采集协议
## 介绍
    OPCUA采集协议是基于OPCUA协议的一种采集协议，主要用于采集OPCUA协议的数据，目前支持的功能有：
        - 读取数据-nodeId方式
        - 认证模式：用户名密码认证，匿名
        - 消息模式 ：支持 encrypt，sign，none
    写数据功能暂未实现，后续会陆续支持。
## 配置
一般会有至少一个自定义协议，关键字段是 `deviceConfig` ，下面给出一个 JSON 样例:

```json
{
  "commonConfig": map[string]interface{}{
  "endpoint":  "opc.tcp://NOAH:53530/OPCUA/SimulationServer",
  "policy":    device.POLICY_BASIC128RSA15,
  "mode":      device.MODE_SIGN,
  "auth":      device.AUTH_ANONYMOUS,
  "username":  "1",
  "password":  "1",
  "timeout":   10,
  "frequency": 500,
  "retryTime": 10,
},
"Opcuanodes": []map[string]interface{}{
{
"tag":         "node1",
"description": "node 1",
"nodeId":      "ns=3;i=1013",
"dataType":    "String",
"value":       "",
},
{
"tag":         "node2",
"description": "node 2",
"nodeId":      "ns=3;i=1001",
"dataType":    "String",
"value":       "",
},
},
}
```

## 字段：

commonConfig 字段：
  - Endpoint：OPC UA服务端的地址，对应gopcua库中的ClientEndpoint
  - Policy：安全策略URL，可以是None、Basic128Rsa15、Basic256、Basic256Sha256中的任意一个，对应gopcua库中的ClientSecurityPolicyURI
  - Mode：安全模式，可以是None、Sign、SignAndEncrypt中的任意一个，对应gopcua库中的ClientSecurityMode
  - Auth：认证模式，可以是Anonymous、UserName、Certificate中的任意一个，对应gopcua库中的ClientAuthMode。
  - Username：用户名，对应gopcua库中的ClientUsername
  - Password：密码，对应gopcua库中的ClientPassword
  - Timeout：超时时间，对应gopcua库中的RequestTimeout
  - Frequency：采集频率，
  - RetryTime：重试次数，

Opcuanodes 点表字段：

- Tag：点表的标签，用于标识点表
-  Description：点表的描述，用于描述点表
- NodeId：点表的NodeId，对应gopcua库中的NodeId
- DataType：点表的数据类型，对应gopcua库中的VariantType
- Value：点表的值，对应gopcua库中的Variant


## 设备数据读取
```lua
		if err := engine.LoadDevice(GENERIC_OPCUA); err != nil {
		t.Fatal(err)
	}
	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{},
		[]string{GENERIC_OPCUA.UUID},
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(data)
			    print(data)
				return true, data
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
		t.Fatal(err)
	}
```
## 设备数据写入
```lua
暂无
```
## 常用函数