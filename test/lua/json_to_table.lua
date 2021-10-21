local json = require "json"

local V1 = json.encode({1, 2, 3, {x = 10}})
print(V1)
local V2 = json.decode('[1,2,3,{"x":10}]')
print(V2[1])
