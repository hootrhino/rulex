-- Copyright (C) 2023 wwhai
--
-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU Affero General Public License as
-- published by the Free Software Foundation, either version 3 of the
-- License, or (at your option) any later version.
--
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU Affero General Public License for more details.
--
-- You should have received a copy of the GNU Affero General Public License
-- along with this program.  If not, see <http://www.gnu.org/licenses/>.
--
-- 主要针对特殊的IEC103协议类设备
--

--[[
    定义一系列设备的报文
--]]
local Devices = {
    {
        Name = "某个设备",
        Encode = function(args)
            -- 你实现
        end,
        Decode = function(args)
            -- 你实现
        end,
    },
}

for __i, Device in ipairs(Devices) do
    local response, err = rulexlib:CtrlDevice("UUID",
        Device.Encode("0001020304AABBCCDD(你的请求报文)"))
    if err ~= nil then
        rulexlib:Log(err)
    else
        local parsedValue = Device.Decode(response)
        print('parsedValue', parsedValue)
    end
end
