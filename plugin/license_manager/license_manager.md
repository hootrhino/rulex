# license manager: 固件证书管理器
固件证书管理器，用来防止盗版或者破解。
## 原理
该插件会向一个服务器发送一个 HTTP 请求，该请求包含了：
1. 本地计算机 MAC 地址
1. 本地软件版本号
1. 本地操作系统
1. 本地硬件架构
2. 或者其他更多

最终会将上述信息用某种算法加密后上传到服务器，服务器验证后返回一个加密数字证书，该证书存放在本地路径下。每次开机启动的时候验证证书即可。

注册伪代码：
```
license_text = request_server("*********")
if check(license_text){
    save(license_text)
}
```

启动认证伪代码:
```
license_text = load_file("/path/app.lic")
if check(license_text){
    start()
}else{
    stop(error-message)
}
```

## 配置
```ini
[plugin.license_manager]
local_addr  = "./rulex/cert.pem"
remote_addr = "http://remote.site/cert/"
net_name    = "eth0"
```

## 接口
远程服务器承担了拉取、验证证书的作用，需要实现对应的接口

```
创建证书
POST remote.site/cert
request:
{
    "mac": "xxx",
    "os": "linux",
    "arch": "amd64",
}
response:
{
    "ret":0,
    "msg": "",
    "data": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----"
}
```

```
验证证书
GET remote.site/cert
request:
{
    "mac": "xxx",
    "data": "xxx" // 使用证书加密mac后的数据，使用base64编码（二进制安全）
}
response:
{
    "ret": 0,
    "msg": "xxx"
}
```