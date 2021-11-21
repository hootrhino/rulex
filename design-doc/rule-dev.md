## 规则引擎使用
### 接口
```go
type Rule struct {
	Id          string      `json:"id"`
	UUID        string      `json:"uuid"`
	Status      RuleStatus  `json:"status"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	VM          *lua.LState `json:"-"`
	From        []string    `json:"from"`
	Actions     string      `json:"actions"`
	Success     string      `json:"success"`
	Failed      string      `json:"failed"`
}

```

### 编写规则回调
```lua

---@diagnostic disable: undefined-global

-- Success
function Success()
    print("======> success")
end
-- Failed
function Failed(error)
    print("======> failed:", error)
end

-- Actions
Actions = {
    function(data)
        return true, data
    end
}

```
### 库函数使用
#### 推送MQTT
```lua
stdlib:DataToMqttServer('id', data)
```
#### 推送Mongo
```lua
stdlib:DataToMongo('id', data)
```
#### 推送HTTP
```lua
stdlib:DataToHttpServer('id', data)
```
#### JSON提取
```lua
Actions = {
	function(data)
	    local V1 = stdlib:JqSelect(".[] | select(.temp > 50000000)", data)
        print("[LUA Actions Callback 1 ===> Data is:", data)
	    print("[LUA Actions Callback 1 ===> .[] | select(.temp >= 50000000)] return => ", stdlib:JqSelect(".[] | select(.temp > 50000000)", data))
		return true, data
	end,
	function(data)
	    local V2 = stdlib:JqSelect(".[] | select(.hum < 20)", data)
	    print("[LUA Actions Callback 2 ===> .[] | select(.hum < 20)] return => ", stdlib:JqSelect(".[] | select(.hum < 20)", data))
		return true, data
	end,
	function(data)
	    local V3 = stdlib:JqSelect(".[] | select(.co2 > 50)", data)
	    print("[LUA Actions Callback 3 ===> .[] | select(.co2 > 50] return => ", stdlib:JqSelect(".[] | select(.co2 > 50)", data))
		return true, data
	end,
	function(data)
	    local V4 = stdlib:JqSelect(".[] | select(.lex > 50)", data)
	    print("[LUA Actions Callback 4 ===> .[] | select(.lex > 50)] return => ", stdlib:JqSelect(".[] | select(.lex > 50)", data))
		return true, data
	end,
	function(data)
		local json = require("json")
		print("[LUA Actions Callback 5, json.decode] ==>",json.decode(data))
		print("[LUA Actions Callback 5, json.encode] ==>",json.encode(json.decode(data)))
		return true, data
	end
}
```
#### 二进制位匹配
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
        print("stdlib:MatchBinary: ", json.encode(result))
        return true, data
    end
}

```
#### 位串转字节
```lua
-- data: 0101010100101001010101010010101010101
stdlib:BitStringToBytes(data)
```
#### 字节转整形
```lua
stdlib:ByteToInt(bytes)
```
#### 取一个字节某个位
```lua
stdlib:GetABitOnByte(index)
```
#### 字节转位串
```lua
stdlib:ByteToBitString(bytes)
```