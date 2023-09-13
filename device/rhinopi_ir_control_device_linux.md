# RhinoPI 红外接收器
RhinoPI自带一路红外接收器，当用遥控器的时候可以触发事件。
## 数据格式
```json
{
    "time":{
        "sec":1694615990,
        "usec":755130
    },
    "type":4,
    "code":4,
    "value":12
}
```

## 脚本示例
```lua
Actions =
{
    function(data)
        rulexlib:Debug(data)
        return true, data
    end
}
```