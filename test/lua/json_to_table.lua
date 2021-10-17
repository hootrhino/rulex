local json = require "json"

local t = {}

local encode = json:encode(t)
print(encode)
local decode = json:decode(encode)
print(decode)
