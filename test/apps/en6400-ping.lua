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
    while true do
        local _, Error = network:Ping("114.114.114.114");
        if Error ~= nil then
            for i = 1, 5, 1 do
                en6400:Led1On();
                time:Sleep(50);
                en6400:Led1Off();
                time:Sleep(50);
            end;
        else
            en6400:Led1On();
            time:Sleep(50);
            en6400:Led1Off();
            time:Sleep(50);
        end;
        time:Sleep(5000);
    end;
    return 0;
end;
