# 双向通信的实现
涉及到开关控制类的，此时其实是网关直接和开关通信。对于这种场景，RULEX也提供了直接支持。
## 下行函数
- `rulex:WriteInStream('id', data)`
    ```lua
    Actions = {
        function(data)
            rulex:WriteInStream('id', data)
            return true, data
        end
    }
    ```
    id: 指的是某个驱动的ID，例如可能是个Modbus客户端。

有时候我们也需要驱动给网关做反馈，用到的函数是：
## 上行函数
-  `rulex:WriteOutStream('id', data)`
    ```lua
    Actions = {
        function(data)
            rulex:WriteOutStream('id', data)
            return true, data
        end
    }
    ```
## 案例
某个产品是16路继电器，其中继电器上的主控MCU每一段时间就会给RULEX上报状态，数据到达RULEX的时候，可用 `rulex:WriteOutStream('id', data)` 来处理转发等