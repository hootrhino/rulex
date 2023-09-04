<!--
 Copyright (C) 2023 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

# RULEX 插件开发指南
## 插件概述
插件是为了给 RULEX 增加扩展没有的功能, 或者外挂用户自己开发的一些服务, 比如你搞了个TCP Server 可以外挂进来，RULEX会做资源管理。插件多用来增加一些和 RULEX 主体功能无关的额外功能。

## 插件接口
下面是插件的接口定义：
```go
type XPlugin interface {
	Init(*ini.Section) error
	Start(RuleX) error
	Service(ServiceArg) ServiceResult
	Stop() error
	PluginMetaInfo() XPluginMetaInfo
}
```

接口解释:
- `Init(*ini.Section) error`
    初始化插件的时候使用，参数为rulex.conf里的参数映射。
- `Start(RuleX) error`
    插件启动接口，参数为 RULEX 接口，可以注入 RULEX 实例，通常相当于一个独立程序的 Main 函数。
- `Service(ServiceArg) ServiceResult`
    插件服务接口，用来向外界提供功能，外界可以通过 RULEX 的 HTTP API 接口来调用这个服务的功能。
- `Stop() error`
    插件停止，这里用来释放资源。
- `PluginMetaInfo() XPluginMetaInfo`
    返回插件的一些原始信息，例如作者等。
## 案例详解

首先我们来看一个最简单的插件示例。
```go
package demo_plugin

import (
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/ini.v1"
)

type DemoPlugin struct {
	uuid string
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{
		uuid: "DEMO01",
	}
}

func (dm *DemoPlugin) Init(config *ini.Section) error {
	return nil
}

func (dm *DemoPlugin) Start(typex.RuleX) error {
	return nil
}
func (dm *DemoPlugin) Stop() error {
	return nil
}

func (hh *DemoPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "DemoPlugin",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

func (cs *DemoPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

```

上面这段代码就是一个空插件，其没有任何实际功能，但是经过加载后会在界面上显示出来。接下来我们重点看一下这几个地方。
1. UUID
    UUID 是一个插件的唯一识别码，这里需要设置好规则，为了给开发者最大的自由，此处 UUID 直接暴露出来，所见即所得，也就是说，你写出来的插件里面的 UUID 是什么，到时候呈现到界面上就是什么。这个很关键，因为前端就是通过 UUID 来控制不同的插件的功能的，如果UUID 冲突或者没有选择合适会影响前端判断。通常 UUID 一旦定了就不会再变，最好用一个有意义的 UUID，例如 RULEX 的API Server 的 UUID 是："HTTP-API-SERVER"。
2. Service 接口
    Service 接口是插件向外界提供功能的入口，外界通过 HTTP 请求进来的参数，会被传到 `Service` 接口的 `arg` 参数。 其中 `ServiceArg`参数的形式如下：
    ```go
    type ServiceArg struct {
        UUID string      `json:"uuid"`
        Name string      `json:"name"`
        Args interface{} `json:"args"`
    }
    type ServiceResult struct {
        Out interface{} `json:"out"`
    }
    ```
    当你需要开发和前端交互的功能的时候，就需要支持该接口，比如请看下面这个功能，也就是RULEX自带的网络测速器, 其中网络测速器的UUID是`ICMPSender`。网络测速器接受一个 IP 列表参数，其实就类似我们操作系统的 `ping 192.168.1.1` 命令，操作形式很简单，对前端而言，构造的JSON如下：
    ```json
    {
        "uuid": "ICMPSender",
        "name": "ping",
        "args": [
            "192.168.1.1"
        ]
    }
    ```
    其中请求会被RULEX处理后转发到 `Service` 接口。`name` 对应 `ServiceArg.Name`，`args` 对应 `ServiceArg.Args`。`Out` 则是返回给前端的值。
    下面这个就是网络测速器的 service 接口部分代码：
    ```go
    func (icmp *ICMPSender) Service(arg typex.ServiceArg) typex.ServiceResult {
        // ping 8.8.8.8
        Fields := logrus.Fields{
            "topic": "plugin/ICMPSenderPing/ICMPSender",
        }
        if arg.Name == "ping" {
            if icmp.pinging {
                glogger.GLogger.WithFields(Fields).Info("ICMPSender pinging now:", arg.Args)
                return typex.ServiceResult{Out: []map[string]interface{}{}}
            }
        }
        // .... 太长了这里不全部展示，请直接看源码
        return typex.ServiceResult{Out: []map[string]interface{}{}}
    }
    ```
    其中有这个日志打印行：
    ```go
    Fields := logrus.Fields{
        "topic": "plugin/ICMPSenderPing/ICMPSender",
    }
    glogger.GLogger.WithFields(Fields).Info("ICMPSender pinging now:", arg.Args)
    ```
    通过这个特殊的 Fields 可以让前端选择日志打印形式。
## 请求响应
Service接口的请求结构如下：
```json
{
    "uuid": "MODBUS_SCANNER",
    "name": "stop",
    "args": {}
}
```
其中`args`是泛型参数，其值类型可能是`string | object | array`。
Service接口的返回结构如下：
```json
{
    "code": 200,
    "msg": "Success",
    "data": {}
}
```
其中`data`是泛型参数，其值类型可能是`string | object | array`。

> 注意：一般泛型类型只和具体插件有关，也就是“到底是什么类型只有插件自己知道”。

## 注册插件
目前注册插件是硬编码形式，是故意这么设计的，具体注册看这个函数即可：https://vscode.dev/github/hootrhino/rulex/blob/dev-v0.6.2/engine/runner.go#L137。
为什么这么做？回顾一下插件定义：
> 插件是为了给 RULEX 增加扩展没有的功能, 或者外挂用户自己开发的一些服务, 比如你搞了个TCP Server 可以外挂进来，RULEX会做资源管理。插件多用来增加一些和 RULEX 主体功能无关的额外功能。

也就是说插件和RULEX是一体的，理应该在编译时就被确认，**所以属于高级开发阶段的东西，而非用户阶段的**。到这里就明白了了，其实类似于Linux内核某些设计一样。

## 总结
插件开发其实很简单，只需搞清楚 `uuid`、`service`、`日志输出` 这三个关键点即可。