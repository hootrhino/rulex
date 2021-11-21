# Rulex 和云端的交互接口
## MQTT Topic 规范
### 状态上报
- Topic: `rulex/status`
- Payload
    ```json
    {
        "clientId":"aabbccdd",
        "state":"running"
    }
    ```
### 拓扑上报
- Topic: `rulex/topology`
- Payload
    ```json
    {
        "clientId":"aabbccdd",
        "network":[
            {
                "name":"modbus",
                "state":"running",
                "type":"direct",
            },
            {
                "name":"snmp",
                "state":"running",
                "type":"direct",
            }
        ]
    }
    ```
### 远程配置

- Topic: `rulex/{clientId}/remoteConfig` ,在线远程下发配置信息到网关里面, 其实就是 rulex 的配置文件的JSON格式:
```ini
.....
name=rulex
log_level=all
log_path=rulex-log.txt
max_queue_size=204800
resource_restart_interval=5000
gomax_procs=2
enable_pprof=false
.....
```
- Payload
    ```json
    {
        //.....
    }
    ```
### 远程操作
- Topic: `rulex/{clientId}/remoteOperate` ,远程创建资源，规则等.
- Payload
    ```json
    {
        "type":"createInEnd",
        "config":{
            "host": "127.0.0.1",
            "port": 2883,
            "topic": "rulex-client-topic-1",
            "clientId": "rulex1",
            "username": "rulex1",
            "password": "******"
        }
    }
    ```
### 日志操作
- Topic: `rulex/{clientId}/logs` ,云端拉取本地日志.
- Payload
    ```json
    {
        "beginTime": "2021-11-21",
        "endTime": "2021-11-22",
        "count": 100,
    }
    ```