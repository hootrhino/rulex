#!/bin/bash

install(){
    local source_dir="$PWD"
    local service_file="/etc/systemd/system/rulex.service"
    local executable="/usr/local/rulex"
    local working_directory="/usr/local/"
    local config_file="/usr/local/rulex.ini"
    local db_file="/usr/local/rulex.db"
cat > "$service_file" << EOL
[Unit]
Description=Rulex Daemon
After=network.target

[Service]
Environment="ARCHSUPPORT=EEKITH3"
WorkingDirectory=$working_directory
ExecStart=$executable run -config=$config_file -db=$db_file
Restart=always
User=root
Group=root
StartLimitInterval=0
RestartSec=5
[Install]
WantedBy=multi-user.target
EOL
    chmod +x $source_dir/rulex
    cp "$source_dir/rulex" "$executable"
    cp "$source_dir/rulex.ini" "$config_file"
    systemctl daemon-reload
    systemctl enable rulex
    systemctl start rulex
    if [ $? -eq 0 ]; then
        echo "Rulex service has been created and extracted."
    else
        echo "Failed to create the Rulex service or extract files."
    fi
    exit 0
}

start(){
    systemctl daemon-reload
    systemctl start rulex
    echo "RULEX started as a daemon."
    exit 0
}
status(){
    systemctl status rulex
}
restart(){
    systemctl stop rulex
    start
}

stop(){
    systemctl stop rulex
    echo "Service Rulex has been stopped."
}

uninstall(){
    local working_directory="/usr/local"
    systemctl stop rulex
    systemctl disable rulex
    rm /etc/systemd/system/rulex.service
    rm $working_directory/rulex
    rm $working_directory/rulex.ini
    rm $working_directory/rulex.db
    rm $working_directory/*.txt
    rm $working_directory/*.txt.gz
    systemctl daemon-reload
    systemctl reset-failed
    echo "Rulex has been uninstalled."
}

check_deps() {
    # Check if ffmpeg is installed
    if command -v ffmpeg >/dev/null 2>&1; then
        echo "ffmpeg is installed"
    else
        echo "Please install ffmpeg"
        exit 1
    fi
    # Check if ttyd is installed
    if command -v ttyd >/dev/null 2>&1; then
        echo "ttyd is installed"
    else
        echo "Please install ttyd"
        exit 1
    fi
    echo "All necessary software is installed"
}
# create a default user
create_user(){
    response=$(curl -X POST -H "Content-Type: application/json" -d '{
    "role": "admin",
    "username": "admin",
    "password": "admin",
    "description": "system admin"
    }' http://127.0.0.1:2580/api/v1/users -w "%{http_code}")

    if [ "$response" = "201" ]; then
        echo "User created"
    else
        echo "User creation failed"
    fi

}
#
#
#
main(){
    case "$1" in
        "install" | "start" | "restart" | "stop" | "uninstall" | "create_user" | "status")
            $1
        ;;
        *)
            echo "Invalid command: $1"
            echo "Usage: $0 <install|start|restart|stop|uninstall|status>"
            exit 1
        ;;
    esac
    exit 0
}
#===========================================
# main
#===========================================
if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root"
    exit 1
else
    check_deps
    main $1
fi