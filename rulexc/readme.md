# RULEXC 客户端工具使用文档
## RULEXC 简介
RULEXC 是 RULEX Client 的意思，是命令行下的客户端工具，帮助我们调试和管理RULEXC。
## 使用文档

- 查看系统参数
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'

```
- 查看入口列表
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 查看单个入口
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 查看出口列表
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 查看单个出口
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 查看单个规则
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 查看单个插件
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 创建入口
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
- 创建出口
```sh
go run ./rulexc.go inend-create --config  '{"name":"test","type":"MQTT","config":{"server":"127.0.0.1","port":1883,"username":"test","password":"test","clientId":"test"},"description":"Description"}'
```
如果需要格式化JSON的工具，可以上:www.json.cn 测试。