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

-- 003e 004c 00cc 00cd
function Main(arg)
    local HexS = "3e4ccccd"
    Debug("Bin2F32Big:" .. HexS .. "->" .. binary:Bin2F32Big(HexS))
    -- 8.9
    Debug("Bin2F32Little:" .. HexS .. "->" .. binary:Bin2F32Little(HexS))
    return 0
end
