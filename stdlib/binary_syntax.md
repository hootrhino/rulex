# 二进制匹配语法
## 规则
一个 `<` 或者 `>` 开头，表示大小端模式，后面紧跟着`K:Length` 键值对,互相之间用空格隔开, `length` 长度是位长.
语法：
```
<k1:length k2:length k2:length ....
```
Demo：
```
<a:16 b:16 c:16 d1:16
```
## 说明
1. 不遵循格式规范回提取失败
2. Key最长建议不要超过32的字符，否则会报错
3. 建议不要匹配太多字节，会影响效率

## 案例
下面看这个读取modbus数据的案例：
```lua
-- Success
function Success()
end
-- Failed
function Failed(error)
    print("Error:", error)
end

-- Actions
Actions =
    {
    --        ┌───────────────────────────────────────────────┐
    -- data = |00 01 00 01|00 01 00 01|00 00 00 01|00 00 00 01|
    --        └───────────────────────────────────────────────┘
    function(data)
        local json = require("json")
        local tb = stdlib:MatchBinary("<a:16 b:16 c:16 d1:16", data, false)
        local result = {}
        result['a'] = stdlib:ByteToInt64(1, stdlib:BitStringToBytes(tb["a"]))
        result['b'] = stdlib:ByteToInt64(1, stdlib:BitStringToBytes(tb["b"]))
        result['c'] = stdlib:ByteToInt64(1, stdlib:BitStringToBytes(tb["c"]))
        result['d1'] = stdlib:ByteToInt64(1, stdlib:BitStringToBytes(tb["d1"]))
        print("stdlib:MatchBinary 2:", json.encode(result))
        return true, data
    end
}

```