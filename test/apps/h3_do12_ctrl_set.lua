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
    while true do
        rhinopi:H3DO1Set(0)
        time:Sleep(1000)
        rhinopi:H3DO1Set(1)
        time:Sleep(1000)
        --
        rhinopi:H3DO2Set(1)
        time:Sleep(1000)
        rhinopi:H3DO2Set(0)
        time:Sleep(1000)
    end
    return 0
end
