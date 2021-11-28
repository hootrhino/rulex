# Rulex 和私有云服务的交互接口
Rulex 作为一个公共组件，***不具备为任何私有云平台或者系统定制的功能***，但是我们可以通过一些资源或者插件来实现和私有平台交互的功能。
目前可以通过 MQTT 协议实现和远程服务器之间的交互。其中如果需要监听远程服务器的消息，首先要创建一个 MQTT 出口, 服务器端建议使用 EMQX 作为代理，关键配置如下:
```json
{
        "host": "127.0.0.1",
        "port": 1883,
        "s2cTopic": "rulex-client-1",
        "topologyTopic": "rulex-topology-1",
        "dataTopic": "rulex-data-1",
        "stateTopic": "rulex-state-1",
        "clientId": "rulex-1",
        "username": "rulex-1",
        "password": "******"
}
```
- `s2cTopic`: 来自服务器的数据
- `topologyTopic`: 拓扑结构上报
- `dataTopic`: 上报自己规则引擎的数据
- `stateTopic`: 上报状态

### RULEX 和私有云交互 Topic 规范
整体架构设计
```
   +-----------------+
   |   EMQX Cluster  |
   +--------^--------+
            |
   +--------+--------+
   |    Gateway      |
   +--^-----------^--+
      |           |
      |           |
   +--+--+     +--+--+
   | D1  |     | D2  |
   +-----+     +-----+
```
下面是Topic规范，注意，`.` 并不是 MQTT 协议规范，这里是为了区分业务的一种表示形式，不要被误导。


| 功能                               | 路径                                   | QoS | 行为      |
| ---------------------------------- | -------------------------------------- | --- | --------- |
| 上报日志                           | upstream.gateway.logs/${client-id}     | 0   | publish   |
| 上报拓扑                           | upstream.gateway.topology/${client-id} | 0   | publish   |
| 上报指令执行结果以及目标节点的状态 | upstream.gateway.state/${client-id}    | 0   | publish   |
| 接受远程消息                       | downstream.gateway.s2c/${client-id}    | 2   | subscribe |
| 规则引擎数据                       | upstream.gateway.publish/${client-id}  | 2   | publish   |
| 设备离线                           | upstream.gateway.disconnected          | 2   | publish   |
| 设备上线                           | upstream.gateway.connected             | 2   | publish   |

***上面的 topic 不是写死的，只是为了配合 EMQX 的推荐值，如果有个性化需求可以自行调整.***

### 消息模板
消息体必须是个JSON，必须包含 `uuid`:
  ```json
  {
      "uuid": "uuid0010101010"
      // other ...
  }
  ```

- 上报日志
  ```json
  {
      "uuid":1,
      "logs":[
          "........"
      ]
  }
  ```
- 上报拓扑
  ```json
  {
      "uuid":1,
      "topology":[
          {"node":"modbus meter1", "state":"running"},
          {"node":"modbus meter2", "state":"running"},
          {"node":"modbus meter3", "state":"running"},
      ]
  }
  ```
  
- 上报指令执行结果以及目标节点的状态
  该功能主要是为了同步设备的状态，比如给某个开关下发了开指令:
  ```json
     {
       "cmdId": "00001",
       "cmd" :"open",
       "sw": [1, 2]
     }
  ```
  此时命令执行完后会有成功或者失败的结果反馈上去, mqtt topic为: `upstream.gateway.state/${client-id}`, 服务端订阅这个 `Topic` 后，可根据 `type` 字段判断类型:
  ```lua
      -- 执行成功
      rulex:finishCmd(CmdId)
  ```
  ```json
     {
       "type": "finishCmd",
       "cmdId" :"00001"
     }
  ```
  
  ```lua
      -- 执行失败
      rulex:failedCmd(CmdId)
  ```
  ```json
     {
       "type": "failedCmd",
       "cmdId" :"00001"
     }
  ```
  
- 接受来自服务端的远程消息格式:

  ```json
  {
      "type": "remoteCmd",
      "cmd": "cmdXXX",
      "args": [
          "AAA"
      ]
  }
  ```
  cmd:
  - `get-state` :通知上报状态
  - `get-topology` :通知上报拓扑
  - `get-log` :通知上报日志

## 开灯Demo
下面以一个Demo来演示：
### LUA 回调
```lua
---
--- 这里展示一个远程发送指令后响应的Demo
--- 假设远程指令是打开开关，然后同步状态到云端,
--- 指令体：{
---            "cmdId" : "hu008987y",
---            "type" : "OPEN",
---            "sn": [
---                   "SN0001",
---                   "SN0002"
---                  ]
---        }
--- 表示打开 SN0001 SN0002 两个开关
---
Actions = {
    function(data)
        local json = require("json")
        local Tb = json.decode(data)
        local CmdId = Tb["cmdId"]
        local Type = Tb["type"]
        local SN = Tb["sn"]
        if Type == "OPEN" then
            local ok = rulex:WriteOutStream('#ID', json.encode({0x00, SN}))
            if ok then
                rulex:finishCmd(CmdId)
            else
                -- 其实没必要显式调用失败，因为服务端超时后就自己直接失败了
                rulex:failedCmd(CmdId)
            end
        end
        if Type == "OFF" then
            local ok = rulex:WriteOutStream('#ID', json.encode({0x01, SN}))
            if ok then
                rulex:finishCmd(CmdId)
            else
                -- 其实没必要显式调用失败，因为服务端超时后就自己直接失败了
                rulex:failedCmd(CmdId)
            end
        end
        return true, data
    end
}

```

### 服务端（go）
```go
package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

//
// 起一个测试环境
// docker run -it --link test-redis:redis --rm redis redis-cli -h redis -p 6379
//
/*
*
* 测试开关打开或者关闭后状态同步机制
*
 */
func Test_Open_Switch(t *testing.T) {
	var redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	defer redisClient.Close()
	//
	//
	//
	requestId := "request-id-001"
	sendCmd(ctx, redisClient, requestId)
	waitResult(ctx, redisClient, requestId)
	time.Sleep(5 * time.Second)

}

/*
*
* 发送指令:当指令下发后马上给redis保存一个指令id，用于等待后期同步
 */
func sendCmd(ctx context.Context, redisClient *redis.Client, requestId string) {
	fmt.Println("Send open cmd to rulex")
	redisClient.Set(ctx, requestId, 0, 5*time.Second)
}

/*
*
* 等待执行结果
*
 */
func waitResult(ctx context.Context, redisClient *redis.Client, requestId string) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				{
					failedCmd(ctx, redisClient, requestId)
					return
				}
			default:
				{
					s := redisClient.Get(ctx, requestId)
					if s.Err() != nil && s.Val() != "" {
						if ok, _ := s.Bool(); ok {
							finishCmd(ctx, redisClient, requestId)
						}
					}
				}
			}
		}
	}(ctx)

}

/*
*
*监听rulex的反馈，如果  rulex:finishCmd(CmdId) 调用了 这里就把redis的值更新
*
 */
func finishCmd(ctx context.Context, redisClient *redis.Client, requestId string) {
	println("finished:" + requestId)

}

/*
*
*
*
 */
func failedCmd(ctx context.Context, redisClient *redis.Client, requestId string) {
	println("failed:" + requestId)
}

```