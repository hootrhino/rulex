# 设备读写操作
设备定义了 Read、Write两个接口，实则定义了设备的对外接口。
## LUA 模板
### 函数原型
- 读: ```rulexlib:ReadDevice(ID) -> data, err```
- 写: ```rulexlib:WriteDevice(ID, []byte{}) -> data, err```
### 写指令到设备：
```lua
    local r1, e1 = rulexlib:WriteDevice("device-uuid", 0, "data")
    if (e1 ~= nil) then
        print('error:', err)
        return false, data
    end
```
### 从设备读数据：
```lua
    local r1, e1 = rulexlib:ReadDevice("device-uuid", 0)
    if (e1 ~= nil) then
        print('error:', err)
        return false, data
    end
```