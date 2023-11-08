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

AppNAME = "每1000ms交替式控制DO"
AppVERSION = "1.0.0"
AppDESCRIPTION = ""
--
-- Main
--

function Main(arg)
    local gpio = { 6, 7 }
    while true do
        for _, value in ipairs(gpio) do
            rhinopi:GPIOSet(value, 0)
            time:Sleep(1000)
            rhinopi:GPIOSet(value, 1)
            time:Sleep(1000)
        end
    end
    return 0
end
