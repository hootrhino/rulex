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

Actions =
{
    function(args)
        local dataT, err = json:J2T(args)
        if (err ~= nil) then
            stdlib:Debug('parse json error:' .. err)
            return true, args
        end
        for key, value in pairs(dataT) do
            local MatchHexS = hex:MatchUInt("hum:[0,1];temp:[2,3]", value['value'])
            local ts = time:Time()
            local Json = json:T2J(
                {
                    tag = key,
                    ts = ts,
                    hum = math:TFloat(MatchHexS['hum'] * 0.1, 2),
                    temp = math:TFloat(MatchHexS['temp'] * 0.1, 2),
                }
            )
            stdlib:Debug(Json)
        end
        return true, args
    end
}
