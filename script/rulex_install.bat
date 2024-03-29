rem Copyright (C) 2023 wwhai
rem
rem This program is free software: you can redistribute it and/or modify
rem it under the terms of the GNU Affero General Public License as
rem published by the Free Software Foundation, either version 3 of the
rem License, or (at your option) any later version.
rem
rem This program is distributed in the hope that it will be useful,
rem but WITHOUT ANY WARRANTY; without even the implied warranty of
rem MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
rem GNU Affero General Public License for more details.
rem
rem You should have received a copy of the GNU Affero General Public License
rem along with this program.  If not, see <http://www.gnu.org/licenses/>.

@echo off

set SERVICE_NAME=RulexService
set BINARY_PATH=%CD%\rulex.exe
set CONFIG_PATH=%CD%\rulex.ini

sc create %SERVICE_NAME% binPath= "%BINARY_PATH% run -config %CONFIG_PATH%" start= auto DisplayName= RulexService

echo Service created successfully.

