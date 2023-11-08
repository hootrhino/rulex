---@diagnostic disable: undefined-global
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

AppNAME = "当检测到DI信号后往tcp服务器发送文本数据"
AppVERSION = "1.0.0"
AppDESCRIPTION = ""
--
-- Main
--

function Main(arg)
    local s = 0
    while true do
        local v, err = rhinopi:GPIOGet(10)
        if err ~= nil then
            break
        else
            if v ~= s then
                local err0 = applib:DataToUdp('udpServerUUID', 'hello gpio10:' .. v)
                print('DataToUdp success? =>', err0 == nil)
            end
            s = v
        end
        time:Sleep(50)
    end
    return 0
end
