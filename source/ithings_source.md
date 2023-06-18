# Ithings 平台支持
## 简介
iThings是一个基于golang开发的轻量级云原生微服务物联网平台.

## 脚本示例
```lua
Actions = {
    function(data)
        print('Data From Ithings:', data)
        local Json = rulexlib:T2J(
            {
                method = 'report',
                params = {
                    tag = 'key',
                    temp = 0.1,
                    hum = 0.1,
                }
            }
        )
        print('Data to Ithings:', Json)
        rulexlib:DataToIthings('INe1e769cff3b9467394564ca78c7bc93b', Json)
        return true, data
    end
}

```
## 注意
当前版本簪不支持网关模式，目前是直连。