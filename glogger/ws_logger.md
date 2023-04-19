# WsLogger
## 简介
这个是针对WEB做的实时日志展示器，其原理是用Websocket推送日志到浏览器。
## 使用
首先web端需要连接到Websocket地址：`127.0.0.1:2580/ws`，连接成功以后必须在5秒内发送固定字符串`WsTerminal`过来，才能连接成功，否则会被强制断开。
## 数据
数据是Json格式，下面是示例日志:
```json
{
    "appId": "rulex",
    "appName": "rulex",
    "file": "rulex/appstack/appstack_runtime.go:77",
    "func": "github.com/i4de/rulex/appstack.(*AppStack).StartApp.func1.1",
    "level": "debug",
    "msg": "App exit",
    "time": "2023-04-19T14:26:54+08:00"
}
```