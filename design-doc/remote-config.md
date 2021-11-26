# Rulex 和私有云服务的交互接口
Rulex 作为一个公共组件，***不具备为任何私有云平台或者系统定制的功能***，但是我们可以通过一些资源或者插件来实现和私有平台交互的功能。
目前可以通过 MQTT 协议实现和远程服务器之间的交互。其中如果需要监听远程服务器的消息，首先要创建一个 MQTT 出口, 服务器端建议使用 EMQX 作为代理，关键配置如下:
```json
{
        "host": "127.0.0.1",
        "port": 1883,
        "s2cTopic": "rulex-client-1",
        "toplogyTopic": "rulex-toplogy-1",
        "dataTopic": "rulex-data-1",
        "stateTopic": "rulex-state-1",
        "clientId": "rulex-1",
        "username": "rulex-1",
        "password": "******"
}
```
- `s2cTopic`: 来自服务器的数据
- `toplogyTopic`: 拓扑结构上报
- `dataTopic`: 上报自己规则引擎的数据
- `stateTopic`: 上报状态

### RULEX 和私有云交互 Topic 规范

| 功能         | 路径                             | QoS | 行为    |
| ------------ | -------------------------------- | --- | ------- |
| 上报日志     | emqx.stream.gateway.logs         | 0   | publish |
| 上报拓扑     | emqx.stream.gateway.toplogy      | 0   | publish |
| 上报自身状态 | emqx.stream.gateway.state        | 0   | publish |
| 接受远程消息 | emqx.stream.gateway.s2c          | 2   | publish |
| 规则引擎数据 | emqx.stream.devices.publish      | 2   | publish |
| 设备离线     | emqx.stream.devices.disconnected | 2   | publish |
| 设备上线     | emqx.stream.devices.connected    | 2   | publish |

***上面的 topic 不是写死的，只是为了配合 EMQX 的推荐值，如果有个性化需求可以自行调整.***

### 消息模板
消息体必须是个JSON，必须包含 `uuid`:
  ```json
  {
      "uuid": "uuid0010101010"
      // other ...
  }
  ```

- 上报日志
  ```json
  {
      "uuid":1,
      "logs":[
          "........"
      ]
  }
  ```
- 上报拓扑
  ```json
  {
      "uuid":1,
      "toplogy":[
          {"node":"modbus meter1", "state":"running"},
          {"node":"modbus meter2", "state":"running"},
          {"node":"modbus meter3", "state":"running"},
      ]
  }
  ```
  
- 上报自身状态
  ```json
    {
        "uuid":1,
        "state":{
            "alloc":12,
            "cpuPercent":[
                0
            ],
            "diskInfo":86,
            "osArch":"windows-amd64",
            "system":31,
            "total":14,
            "version":"0.0.0-4b22a5e74f32bdc"
        }
    }
  ```
  
- 接受远程消息

  ```json
  {
      "cmd":"cmd",
      "args":[
          "........"
      ]
  }
  ```
  cmd:
  - `get-state` :通知上报状态
  - `get-toplogy` :通知上报拓扑
  - `get-log` :通知上报日志