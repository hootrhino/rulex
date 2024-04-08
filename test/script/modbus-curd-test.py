# Copyright (C) 2024 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

import requests
import json

url = "http://127.0.0.1:2580/api/v1/devices/create"

payload = json.dumps({
  "name": "Modbus设备测试",
  "type": "GENERIC_MODBUS",
  "description": "Modbus设备测试",
  "gid": "DROOT",
  "config": {
    "commonConfig": {
      "mode": "UART",
      "autoRequest": True,
      "enableOptimize": False
    },
    "portUuid": "COM12"
  },
  "schemaId": ""
})
headers = {
  'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
