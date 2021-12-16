# 二进制匹配语法
## 规则
一个 `<` 或者 `>` 开头, 表示大小端模式, 后面紧跟着`K:Length` 键值对,互相之间用空格隔开, `length` 长度是位长.
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
2. Key最长建议不要超过32的字符, 否则会报错
3. 建议不要匹配太多字节, 会影响效率

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
        local tb = rulex:MatchBinary("<a:16 b:16 c:16 d1:16", data, false)
        local result = {}
        result['a'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["a"]))
        result['b'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["b"]))
        result['c'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["c"]))
        result['d1'] = rulex:ByteToInt64(1, rulex:BitStringToBytes(tb["d1"]))
        print("rulex:MatchBinary test:", json.encode(result))
        return true, data
    end
}

```
## 其他实现

除了用RulexAPI, 另外也可以用lua的原生语法, 但是可能和 LUA 版本有关, 下面的函数都是 `LUA 5.4.1` 测试过的, 但是在 `LUA 5.1.x` 就会有问题:
```lua
-- 首先定义一个字符串
local str = "012abcd"
print("str = " .. str)

-- 使用常规方式
print("\nafter string.byte(str,1,4)")
print(string.byte(str, 1, 4))

-- 使用另一种表现方式
print("\nafter str:byte(1,4)")
print(str:byte(1, 4))

-- 使用负数索引
print("\nafter str:byte(-2,-1)")
print(str:byte(-2, -1))

-- 当参数i大于j时
print("\nafter str:byte(2,1)")
print(str:byte(2, 1))

-- 当索引无效时
print("\nafter str:byte(2000,1000000)")
print(str:byte(2000, 1000000))

-- 字符转换
-- 转换第一个字符
print(string.byte("Lua"))
-- 转换第三个字符
print(string.byte("ABC", 1, 3))
-- 转换末尾第一个字符
print(string.byte("Lua", -1))
-- 第二个字符
print(string.byte("Lua", 2))
-- 转换末尾第二个字符
print(string.byte("Lua", -2))

print(string.unpack(">B", "012", 1))
print(string.unpack("B", str))

```