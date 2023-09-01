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
    function(data)
        local dataT, err = rulexlib:J2T(data)
        if (err ~= nil) then
            applib:Debug('parse json error:' .. err)
            return true, data
        end
        for key, value in pairs(dataT) do
            local MatchHexS = rulexlib:MatchUInt("temp:[0,1];hum:[2,3]", value['value'])
            local ts = rulexlib:TsUnixNano()
            local Json = rulexlib:T2J(
                {
                    tag = key,
                    ts = ts,
                    temp = MatchHexS['temp'] * 0.1,
                    hum = MatchHexS['hum'] * 0.1,
                }
            )
            rulexlib:Debug(Json)
        end
        return true, data
    end
}
