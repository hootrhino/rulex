---@diagnostic disable: undefined-global
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

Actions = {
    function(args)
        Debug(args)
        local JsonT = json:J2T(args)
        local t = {
            id = string:MakeUid(),
            method = "thing.event.property.post",
            params = {
                ais_data = JsonT['ais_data']
            }
        }
        local jsons = json:T2J(t)
        Debug(jsons)
        local error = data:ToMqtt('OUTQAQXBVCU', jsons)
        if error ~= nil then
            Throw(error)
        end
        return true, args
    end
}
