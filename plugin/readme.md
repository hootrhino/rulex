## 概述
插件是为了给 `RULEX` 增加扩展没有的功能, 或者外挂用户自己开发的一些服务, 比如你搞了个TCP Server 可以外挂进来，RULEX会做资源管理。

## 内置
目前内置了两个插件：
1. Http api server
   这是Rulex 的核心API插件，主要用于提供 Http restapi 服务， 默认端口是:`2580`
2. Simple mqtt server
   这是个简单的 Mqtt 服务器，可以用来做测试，实际性能没有测过，不过应该还行，10000以下的设备没啥问题， 默认端口是:`2883`
## 开发
开发的时候，需要把配置加进`rulex.ini`，然后在init接口中可以拿到配置，转换成具体的配置即可:
```go
// 配置结构体定义
type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
}
// ....
// 初始化的时候转换
func (s *MqttServer) Init(config *ini.Section) error {
    var mainConfig _serverConfig
    if err := utils.InIMapToStruct(config, &mainConfig); err != nil {
        return err
    }
    s.Host = mainConfig.Host
    s.Port = mainConfig.Port
    return nil
}

```