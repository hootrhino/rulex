-- Copyright (C) 2023 wwhai
--
-- This program is free software= you can redistribute it and/or modify
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
-- along with this program.  If not, see <http=//www.gnu.org/licenses/>.

function Main(arg)
    local dataTable = {
        device_uuid = 1,
        recv_time = "2023-11-28T11=11=36+08=00",
        bat_voltage = 0,
        longitude = 0,
        latitude = 0,
        air_height = 0,
        water_temp = 19.3171,
        salinity = 754.523,
        dissolved_oxygen = 0,
        ph_value = 5.94821,
        wind_speed = 1.02,
        wind_direction = 12,
        air_temp = 21.1,
        air_pressure = 102.3,
        air_humidity = 69.1,
        noise = 42.9,
        wave_height = 0,
        mean_wave_period = 0,
        peak_wave_period = 0,
        mean_wave_direction = 0
    }
    local JsonString = json:T2J(dataTable)
    local Value = http:Post("http://127.0.0.1:6003/api", JsonString)
    Debug("Http Post:" .. Value)
    return 0
end
