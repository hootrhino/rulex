# MQTT Broker
该插件是一个很简单轻巧的MQTT Broker，可以连接上万规模的子设备，当然实际情况下也不建议连接这么多，控制在100个以内即可，你可以取决于自己的硬件性能请酌情尝试。

## 功能
- 完整支持 MQTT3.1.1，并且带有 QOS 控制
- 基本的设备列表获取

## 配置
```ini
[plugin.mqtt_server]
#
# Enable
#
enable = false
#
# Server host, default allow all
#
host = 0.0.0.0
#
# Server port
#
port = 1883
```
### 参数说明
- enable: 是否开启
- port：监听端口