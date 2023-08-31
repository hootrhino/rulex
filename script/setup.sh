# Copyright (C) 2023 wwhai
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
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#!/bin/bash
#!/bin/bash

service_file="/etc/systemd/system/rulex.service"
zip_file="app.zip"
extracted_folder="extracted_files"
executable="/usr/local/rulex"
config_file="/usr/local/rulex.ini"


if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root"
    exit 1
fi


cat > "$service_file" << EOL
[Unit]
Description=Rulex Daemon
After=network.target

[Service]
ExecStart=$executable run -config=$config_file
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
EOL

# 解压缩
unzip "$zip_file" -d "$extracted_folder"

# 将 rulex 和 rulex.ini 移动到指定位置
mv "$extracted_folder/rulex" /usr/local/
mv "$extracted_folder/rulex.ini" /usr/local/
rm -r "$extracted_folder"

systemctl daemon-reload
systemctl enable rulex
systemctl start rulex

echo "Rulex service unit file and files from app.zip have been created and extracted."
