# 摄像头推流

## 介绍

从本地或者远程摄像头拉流,支持本地USB和远程RTSP，**注意该功能只支持Linux系统**。

## 配置

这是设备启动所需的必要配置, 以OPCUA的为例:

```go
{
	MaxThread  int    `json:"maxThread"`  // 最大连接数, 防止连接过多导致摄像头拉流失败
	InputMode  string `json:"inputMode"`  // 视频输入模式：RTSP | LOCAL
	Device     string `json:"device"`     // 本地视频设备路径，在输入模式=LOCAL时生效
	RtspUrl    string `json:"rtspUrl"`    // 远程视频设备地址，在输入模式=RTSP时生效
	OutputMode string `json:"outputMode"` // 输出模式：JPEG_STREAM | RTSP_STREAM
	OutputAddr string `json:"outputAddr"` // 输出地址, 格式为: "Ip:Port",例如127.0.0.1:7890
}
```

## 维护

- <cnwwhai@gmail.com>
