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


function Main(arg)
    -- MatchHexS 是一个Map结构，K为表达式的name，Value为十六进制字符串，
    -- binary:MatchUInt:表示提取值转换成4字节无符号数，也就是原始字节。
    -- "hum:[0,1];temp:[2,3]": 表示提取第 0 1 个字节，赋值给hum；提取第 2 3个字节赋值给humi
    local MatchHexS = hex:MatchUInt("hum:[0,1];temp:[2,3]", "000102030405060708")
    -- 大端输出温度
    Debug("Bin2F32Big hum: " .. binary:Bin2F32Big(MatchHexS['hum']))
    -- 大端输出湿度
    Debug("Bin2F32Big temp: " .. binary:Bin2F32Big(MatchHexS['temp']))
    -- 小端输出温度
    Debug("Bin2F32Little hum: " .. binary:Bin2F32Little(MatchHexS['hum']))
    -- 小端输出湿度
    Debug("Bin2F32Little temp: " .. binary:Bin2F32Little(MatchHexS['temp']))
    return 0
end
