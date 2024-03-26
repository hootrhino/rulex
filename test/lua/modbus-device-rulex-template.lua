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
--
-- 采集到的数据格式如下:
-- [
--     {
--         "tag": "a15",
--         "alias": "a15",
--         "slaverId": 15,
--         "lastFetchTime": 1708873997979,
--         "value": "1.000000"
--     }
-- ]
Actions = {
    -- args 是JSON字符串
    function(args)
        local JsonT = json:J2T(args)
        for _, dt in ipairs(JsonT) do
            Debug('** tag=' .. dt.tag)
            Debug('** alias=' .. dt.alias)
            Debug('** slaverId=' .. dt.slaverId)
            Debug('** lastFetchTime=' .. dt.lastFetchTime)
            Debug('** value=' .. dt.value)
        end
        return true, args
    end
}
