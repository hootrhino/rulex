# 设备网络发现插件
该插件本意被设计用来做设备互相发现用，但是目前暂时还未进行开发工作。
## 获取节点信息
该插件支持监听局域网扫描，只需向节点指定的 UDP 端口（默认1994）发送"NODE_INFO"即可扫描到节点的信息。
```json
{
    "allocMem":14,
    "cpuPercent":0,
    "diskInfo":69,
    "osArch":"windows-amd64",
    "startedTime":"2023-04-13 22:42:34",
    "systemMem":30,
    "totalMem":16,
    "version":{
        "Version":"v0.4.4",
        "ReleaseTime":"2023-04-13 21:44:33"
    }
}
```

## 节点发现配置
```ini
[plugin.netdiscover]
#
# Enable
#
enable = false
#
# Server host, default allow all
#
listen_host = 0.0.0.0
#
# Server port
#
listen_port = 1994
```