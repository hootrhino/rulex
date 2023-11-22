# 摄像头推流

## 介绍

从本地或者远程摄像头拉流，支持本地USB和远程RTSP，Windows 下环境搭建可参考该文档: https://p.kdocs.cn/s/SHGXWBAA3U 。

## 配置

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
## 测试
勃播放地址:`ws://127.0.0.1:9400/ws?token=WebRtspPlayer&liveId=a97607e47c81d43dba8ef6fa48a2cd45`,其中：
- URL: 固定路径`ws://127.0.0.1:9400/ws`
- token：固定值`WebRtspPlayer`
- liveId：播放源的名称的**md5Hash**,例如`USB2.0 PC CAMERA`的 liveId 是 `a97607e47c81d43dba8ef6fa48a2cd45`。

## 维护

- <cnwwhai@gmail.com>
