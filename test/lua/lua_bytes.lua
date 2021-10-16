-- 首先定义一个字符串
local str = "012abcd"
print("str = "..str)

-- 使用常规方式
print("\nafter string.byte(str,1,4)")
print(string.byte(str,1,4))

-- 使用另一种表现方式
print("\nafter str:byte(1,4)")
print(str:byte(1,4))

-- 使用负数索引
print("\nafter str:byte(-2,-1)")
print(str:byte(-2,-1))

-- 当参数i大于j时
print("\nafter str:byte(2,1)")
print(str:byte(2, 1))

-- 当索引无效时
print("\nafter str:byte(2000,1000000)")
print(str:byte(2000,1000000))