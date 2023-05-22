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
可以通过下面这个HTML页面来测试效果。
```html
<!DOCTYPE html>
<html lang="zh-CN">

<head>
    <meta charset="UTF-8">
    <meta name="viewport"
        content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0,minimal-ui:ios">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
    <link rel="stylesheet" href="">
    <script src=""></script>
</head>

<body>
    <div>
        <img src="http://127.0.0.1:8080" />
        <img src="http://127.0.0.1:8080" />
        <img src="http://127.0.0.1:8080" />
    </div>
    <div>
        <img src="http://127.0.0.1:8080" />
        <img src="http://127.0.0.1:8080" />
        <img src="http://127.0.0.1:8080" />
    </div>
</body>

</html>
```
## 维护

- <cnwwhai@gmail.com>
