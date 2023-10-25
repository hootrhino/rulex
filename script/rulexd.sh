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
RestartSec=2
StartLimitInterval=0
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
    if [ -d "$working_directory" ]; then
        cd "$working_directory" || exit 1
        if [ -e "rulex" ]; then
            rm "rulex"
            echo "Deleted 'rulex' in $working_directory"
        fi

        if [ -e "rulex.ini" ]; then
            rm "rulex.ini"
            echo "Deleted 'rulex.ini' in $working_directory"
        fi

        if [ -e "rulex.db" ]; then
            rm "rulex.db"
            echo "Deleted 'rulex.db' in $working_directory"
        fi

        if [ -n "$(find . -maxdepth 1 -name '*.txt' -print -quit)" ]; then
            rm *.txt
            echo "Deleted .txt files in $working_directory"
        fi

        if [ -n "$(find . -maxdepth 1 -name '*.txt.gz' -print -quit)" ]; then
            rm *.txt.gz
            echo "Deleted .txt.gz files in $working_directory"
        fi

    fi
    systemctl daemon-reload
    systemctl reset-failed
    echo "Rulex has been uninstalled."
}
# create a default user
create_user(){
    # 检查是否提供了足够的参数
    if [ $# -ne 2 ]; then
        echo "Missing username and password, example: create_user user1 1234"
        exit 1
    fi
    param1="$1"
    param2="$2"

    response=$(curl -X POST -H "Content-Type: application/json" -d '{
    "role": "admin",
    "username": ${param1},
    "password": ${param2},
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
    main $1
fi