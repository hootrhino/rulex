# 通用IOTHUB
## 简介
本功能主要是对通用IOTHUB提供了支持。

## 脚本示例
```lua
Actions = {
    function(args)
        print('Data From IotHUB:', data)
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
        print('Data to IotHUB:', Json)
        data:ToIotHUB('INe1e769cff3b9467394564ca78c7bc93b', Json)
        return true, args
    end
}

```