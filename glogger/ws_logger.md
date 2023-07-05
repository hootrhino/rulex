# WsLogger
## 简介
这个是针对WEB做的实时日志展示器，其原理是用Websocket推送日志到浏览器。
## 使用
首先web端需要连接到Websocket地址：`ws://127.0.0.1:2580/ws`，连接成功以后必须在5秒内发送固定字符串`WsTerminal`过来，才能连接成功，否则会被强制断开。

## Websocket Broker
当前日志系统通过Websocket来交互，后端统一把日志打到Websocket客户端，然后客户端来进行筛选其用途。
日志的格式如下：
```json
{
    "appId":"rulex",
    "appName":"rulex",
    "topic": "app/UUID12345678/log",
    "file":"appstack_runtime.go:86",
    "func":"github.com/hootrhino/rulex/appstack.(*AppStack).StartApp.func1",
    "level":"debug",
    "msg":"Ready to run app:APPefeebdf253544730a9dc38e15354d2d4-AAAA-1.0.0",
    "time":"2023-06-29T16:46:00+08:00"
}
```
关键字段：
- level: 日志的级别
- msg: 日志本体
- topic: 用来标识该日志的一个作用, 例如`app/UUID12345678/log`表示`UUID12345678`这个app的运行日志。

### Topic 规范
Topic 具备一定的格式规范，但是仅仅为了区分业务，格式不会对功能造成影响。

格式如下：
$$
/业务名/模块/······
$$

例如网络测速的日志Topic:
$$
/plugin/ICMPSenderPing/UUID
$$

前端只需要知道Topic即可知道其具体的含义。

### 常见Topic

| Ws log topic                     | 用途                 |
| -------------------------------- | -------------------- |
| plugin/ICMPSenderPing/ICMPSender | 网络测速插件的日志   |
| rule/test/$UUID                  | 某个规则测试日志     |
| rule/log/$UUID                   | 某个规则运行时的日志 |
| app/console/$UUID                | 某个轻量应用运行日志 |
### 常见日志示例
1. 测速
    ```json
    {
        "appId":"rulex",
        "level":"info",
        "msg":"[Count:3] Ping Reply From [192.168.1.1]: time=505.8µs ms TTL=128",
        "time":"2023-06-30T13:04:06+08:00",
        "topic":"plugin/ICMPSenderPing/ICMPSender"
    }
    ```
2. Rule 调试 log
    ```json
    {
        "appId":"rulex",
        "file":"C:/Users/wangwenhai/workspace/rulex/plugin/http_server/rule_api.go:580",
        "func":"github.com/hootrhino/rulex/plugin/http_server.TestSourceCallback",
        "level":"debug",
        "msg":"string",
        "time":"2023-06-30T17:52:31+08:00",
        "topic":"rule/test/INa56de94aa22340c89cfab091a53d074f"
    }
    ```
3. Rule 运行时 log
    ```json
    {
        "appId":"rulex",
        "file":"C:/Users/wangwenhai/workspace/rulex/plugin/http_server/rule_api.go:580",
        "func":"github.com/hootrhino/rulex/plugin/http_server.TestSourceCallback",
        "level":"debug",
        "msg":"string",
        "time":"2023-06-30T17:52:31+08:00",
        "topic":"rule/log/INa56de94aa22340c89cfab091a53d074f"
    }
    ```
4. APP 运行输出
    ```json
    {
        "appId":"rulex",
        "file":"C:/Users/wangwenhai/workspace/rulex/rulexlib/log_lib.go:35",
        "func":"github.com/hootrhino/rulex/rulexlib.DebugAPP.func1",
        "level":"debug",
        "msg":"2023-06-30 17:31:35",
        "time":"2023-06-30T17:31:35+08:00",
        "topic":"app/console/APP66c580e7b2c04aa18c30164973ec1d76"
    }
    ```
