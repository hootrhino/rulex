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
-- along with this program.  If not, see <http://www.gnu.org/licenses/>.

-- {
--     "id": "1",
--     "method":"thing.event.property.post",
--     "params": {
--          "temp1":1
--     }
-- }

Actions = {
    function(args)
        local dataT, err = json:J2T(args)
        if (err ~= nil) then
            stdlib:Throw('parse json error:' .. err)
            return true, args
        end
        for _, value in pairs(dataT) do
            local params = {}
            params[value['tag']] = value.value
            local json = json:T2J({
                id = "1",
                method = "thing.event.property.post",
                params = params
            })
            stdlib:Debug(json)
            local err = data:ToMqtt('OUTSKGLIQJX', json)
            if err ~= nil then
                stdlib:Throw(err)
                return true, args
            end
        end
        return true, args
    end
}
