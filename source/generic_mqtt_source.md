# 通用MQTT使用方法

## 简介
通用MQTT南向资源，支持连接指定mqtt服务并订阅多个topic

## 配置
```json
{
    "clientId": "eekit9791", // 客户端ID
    "host": "127.0.0.1", // 服务地址
    "port": 1883, // 端口
    "username": "hootrhino", // 用户名
    "password": "12345678", // 密码
    "qos": 1, // 数据质量0，1，2
    "subTopics": [ // 所需订阅的topic
        "/topic1/#",
        "/topic2/#",
        "/topic3/#"
    ]
}
```


## 规则示例
```lua
-- args示例
-- {"topic":"/topic1/dev","payload":"WyJzdHIxIiwgInN0cjIiXQ=="}

Actions = {
	function(args)
		stdlib:Debug(args)
		local body = json:J2T(args)
        stdlib:Debug(body["topic"])
		local payloadTable, err = json:base64J2T(body["payload"])
		if err ~= nil then
			stdlib:Debug(err)
			return false, args
		end
        stdlib:Debug(payloadTable[1])
        stdlib:Debug(payloadTable[2])
		return true, args
	end
}
```
