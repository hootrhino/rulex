# 摄像头推流

## 介绍

从本地或者远程摄像头拉流，支持本地USB和远程RTSP，Windows 下环境搭建可参考该文档: https://p.kdocs.cn/s/SHGXWBAA3U 。

## 配置
本地摄像头：
```json
{
    "name": "GENERIC_CAMERA",
    "type": "GENERIC_CAMERA",
    "gid": "DROOT",
    "description": "GENERIC_CAMERA",
    "config": {
        "inputMode": "LOCAL",
        "device": "USB2.0 PC CAMERA",
        "rtspUrl": "",
        "OutputEncode": "H264_STREAM"
    }
}
```
远程摄像头：
```json
{
    "name": "GENERIC_CAMERA",
    "type": "GENERIC_CAMERA",
    "gid": "DROOT",
    "description": "GENERIC_CAMERA",
    "config": {
        "inputMode": "RTSP",
        "device": "",
        "rtspUrl": "rtsp://192.168.1.210:554/av0_0",
        "OutputEncode": "H264_STREAM"
    }
}
```
## 测试
使用下面的页面测试：
- https://www.zngg.net/tool/detail/FlvPlayer或者
- https://xqq.im/mpegts.js/demo/

播放地址:`ws://127.0.0.1:9400/ws?token=WebRtspPlayer&liveId=a97607e47c81d43dba8ef6fa48a2cd45`,其中：
- URL: 固定路径`ws://127.0.0.1:9400/ws`
- token：固定值`WebRtspPlayer`
- liveId：播放源的名称的**md5Hash**,例如`USB2.0 PC CAMERA`的 liveId 是 `a97607e47c81d43dba8ef6fa48a2cd45`。

## 维护

- <cnwwhai@gmail.com>
