## BACnet设备

### 设备type值
`GENERIC_BACNET_IP`
### 配置示例
```json
{
    "commonConfig": {
        "ip": "192.168.0.197",
        "port": 47808,
        "localPort": 0,
        "interval": 10
    },
    "nodeConfig": [
        {
            "tag": "t1",
            "type": 0,
            "id": 0
        },
        {
            "tag": "tn",
            "type": 0,
            "id": 0
        }
    ]
}
```

### 配置字段
1. 字段含义
- ip:  bacnet设备ip
- port:  bacnet端口，通常是47808
- localPort:  本地监听端口，填0表示默认47808（有的模拟器必须本地监听47808才能正常交互）
- interval:  采集间隔，单位秒
- nodeConfig:  点位列表
- tag":  tagId，必填
- type":  object类型，必填。（下拉框，枚举值）
- id":  objectId，必填。（范围0-4194303)
2. 特殊枚举
nodeConfig.type字段枚举，**前端展示字母，json中存数值**。
- 0: AI
- 1: AO
- 2: AV
- 3: BI
- 4: BO
- 5: BV
- 13: MSI
- 14: MSO
- 19: MSV

### 数据格式
```json
{
    "tag1": "1",
    "tag2": "2",
    "tag3": "3"
}
```