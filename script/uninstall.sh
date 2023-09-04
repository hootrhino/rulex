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
working_directory="/usr/local"

if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root"
   exit 1
fi

systemctl stop rulex
systemctl disable rulex
rm /etc/systemd/system/rulex.service
rm $working_directory/rulex
rm $working_directory/rulex.ini
rm $working_directory/rulex.db
rm $working_directory/*rulex-log.txt
rm $working_directory/*rulex-lua-log.txt
systemctl daemon-reload
systemctl reset-failed
echo "Rulex uninstalled."