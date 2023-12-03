# 内部消息源
可以使用这个数据源来监控系统内部的消息，比如设备离线，资源连接失败等。

## 消息类型
- SOURCE: 南向事件
- DEVICE: 设备事件
- TARGET: 北向事件
- SYSTEM: 系统内部事件
- HARDWARE: 硬件事件

## 示例
### 设备上线:
```json
{
	"type" :"DEVICE",
	"event":"event.connected",
	"ts":121312431432,
	"device_info":{
		"uuid":"UUID1234567",
		"name":"温湿度计"
	}
}

```

### 设备离线:
```json
{
	"type" :"DEVICE",
	"event":"event.disconnected",
	"ts":121312431432,
	"device_info":{
		"uuid":"UUID1234567",
		"name":"温湿度计"
	}
}

```