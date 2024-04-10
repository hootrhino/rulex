-- Copyright (C) 2024 wwhai
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
-- along with this program.  If not, see <https://www.gnu.org/licenses/>.

function Main(arg)
    local acc1 = 300
    local acc2 = 100
    while true do
        -- DEVICENKRZFRYW
        local R = dataschema:Update('DEVICENKRZFRYW', json:T2J({
            a = acc1,
            b = acc2
        }))
        acc1 = acc1 + 1
        acc2 = acc2 + 1
        Debug("dataschema:Update acc1:" .. acc1 .. "; R=", R)
        Debug("dataschema:Update acc2:" .. acc2 .. "; R=", R)
        time:Sleep(1000)
    end
    return 0
end
