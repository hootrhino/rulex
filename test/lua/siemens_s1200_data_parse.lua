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

-- Actions
-- {
--     "tag":"Value",
--     "type":"DB",
--     "frequency":0,
--     "address":1,
--     "start":100,
--     "size":16,
--     "value":"00000001000000020000000300000004"
-- }
Actions =
{
    function(args)
        local dataT, err = json:J2T(args)
        if (err ~= nil) then
            Debug('parse json error:' .. err)
            return true, args
        end
        for key, value in pairs(dataT) do
            --data: 00000001000000020000000300000004
            local MatchHexS = hex:MatchUInt("a:[0,3];b:[4,7];c:[8,11];d:[12,15]", value['value'])
            local ts = time:Time()
            local Json = json:T2J(
                {
                    tag = key,
                    ts = ts,
                    a = MatchHexS['a'],
                    b = MatchHexS['b'],
                    c = MatchHexS['c'],
                    d = MatchHexS['d'],
                }
            )
            Debug(Json)
        end
        return true, args
    end
}
