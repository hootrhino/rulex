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
