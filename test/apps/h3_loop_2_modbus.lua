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

AppNAME = "双路Modbus50毫秒轮询测试"
AppVERSION = "1.0.0"
AppDESCRIPTION = ""
--
-- Main
--

function Main(arg)
    while true do
        local result, err = applib:CtrlDevice('$UUID', "010300010001D5CA")
        if err ~= nil then
            stdlib:Debug("!!!*** CtrlDevice Error=>" .. err)
        else
            stdlib:Debug("√√√*** CtrlDevice result=>" .. result)
        end
        time:Sleep(50)
    end
end
