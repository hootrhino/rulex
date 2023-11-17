#!/bin/bash
RESET='\033[0m'
RED='\033[31m'
BLUE='\033[34m'
YELLOW='\033[33m'

# 打印红色文本
echo_red() {
    echo -e "${RED}$1${RESET}"
}

# 打印蓝色文本
echo_blue() {
    echo -e "${BLUE}$1${RESET}"
}

# 打印黄色文本
echo_yellow() {
    echo -e "${YELLOW}$1${RESET}"
}
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
    cp "$source_dir/license.*" /usr/local/
    systemctl daemon-reload
    systemctl enable rulex
    systemctl start rulex
    if [ $? -eq 0 ]; then
        echo "[√] Rulex service has been created and extracted."
    else
        echo "[x] Failed to create the Rulex service or extract files."
    fi
    exit 0
}

start(){
    systemctl daemon-reload
    systemctl start rulex
    echo "[√] RULEX started as a daemon."
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
    echo "[√] Service Rulex has been stopped."
}
remove_files(){
    if ls $1 1> /dev/null 2>&1; then
        rm $1
        echo "[!] $1 files removed."
    else
        echo "[#] $1 files not found. No need to remove."
    fi
}
uninstall(){
    local working_directory="/usr/local"
    systemctl stop rulex
    systemctl disable rulex
    remove_files /etc/systemd/system/rulex.service
    remove_files $working_directory/rulex
    remove_files $working_directory/rulex.ini
    remove_files $working_directory/rulex.db
    remove_files $working_directory/*.txt
    remove_files $working_directory/upload/
    remove_files $working_directory/license.*
    remove_files $working_directory/*.txt.gz
    systemctl daemon-reload
    systemctl reset-failed
    echo "[√] Rulex has been uninstalled."
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
        echo "[√] User created"
    else
        echo "[x] User creation failed"
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
            echo "[x] Invalid command: $1"
            echo "[?] Usage: $0 <install|start|restart|stop|uninstall|status>"
            exit 1
        ;;
    esac
    exit 0
}
#===========================================
# main
#===========================================
if [ "$(id -u)" != "0" ]; then
    echo "[x] This script must be run as root"
    exit 1
else
    main $1
fi