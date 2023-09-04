#!/bin/bash
source_dir="$PWD"
#
service_file="/etc/systemd/system/rulex.service"
executable="/usr/local/rulex"
working_directory="/usr/local/"
config_file="/usr/local/rulex.ini"
db_file="/usr/local/rulex.db"

#
if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root"
    exit 1
fi

#
cat > "$service_file" << EOL
[Unit]
Description=Rulex Daemon
After=network.target

[Service]
WorkingDirectory=$working_directory
ExecStart=$executable run -config=$config_file -db=$db_file
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
EOL

#
cp "$source_dir/rulex" "$executable"
cp "$source_dir/rulex.ini" "$config_file"

#
systemctl daemon-reload
systemctl enable rulex
systemctl start rulex

#
if [ $? -eq 0 ]; then
    echo "Rulex service has been created and extracted."
else
    echo "Failed to create the Rulex service or extract files."
fi
