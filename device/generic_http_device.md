## Http 采集器
主要用来请求远程HTTP接口，采集这类设备的数据。

## 配置

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "uuid": "DEVICEP6VNYAAK",
        "gid": "DROOT",
        "name": "HTTP请求设备",
        "type": "GENERIC_HTTP_DEVICE",
        "state": 1,
        "config": {
            "commonConfig": {
                "autoRequest": true,
                "frequency": 1000,
                "timeout": 3000
            },
            "httpConfig": {
                "headers": {
                    "token": "12345"
                },
                "url": "http://127.0.0.1:8080"
            }
        },
        "description": ""
    }
}
```