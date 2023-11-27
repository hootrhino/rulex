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
    local WORKING_DIRECTORY="/usr/local/"
    local config_file="/usr/local/rulex.ini"
    local db_file="/usr/local/rulex.db"
cat > "$service_file" << EOL
[Unit]
Description=Rulex Daemon
After=network.target

[Service]
Environment="ARCHSUPPORT=EEKITH3"
WorkingDirectory=$WORKING_DIRECTORY
ExecStart=$executable run -config=$config_file -db=$db_file
ConditionPathExists=!/var/run/rulex-upgrade.lock
Restart=always
User=root
Group=root
StartLimitInterval=0
RestartSec=5
[Install]
WantedBy=multi-user.target
EOL
    chmod +x $source_dir/rulex
    echo "[.] Copy $source_dir/rulex to $WORKING_DIRECTORY."
    cp "$source_dir/rulex" "$executable"
    echo "[.] Copy $source_dir/rulex.ini to $WORKING_DIRECTORY."
    cp "$source_dir/rulex.ini" "$config_file"
    echo "[.] Copy $source_dir/license.key to /usr/local/license.key."
    cp "$source_dir/license.key" "/usr/local/license.key"
    echo "[.] Copy $source_dir/license.lic to /usr/local/license.lic."
    cp "$source_dir/license.lic" "/usr/local/license.lic"
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
remove_files() {
    if [ -e "$1" ]; then
        if [[ $1 == *"/upload"* ]]; then
            rm -rf "$1"
        else
            rm "$1"
        fi
        echo "[!] $1 files removed."
    else
        echo "[*] $1 files not found. No need to remove."
    fi
}

uninstall(){
    systemctl stop rulex
    systemctl disable rulex
    remove_files /etc/systemd/system/rulex.service
    remove_files $WORKING_DIRECTORY/rulex
    remove_files $WORKING_DIRECTORY/rulex.ini
    remove_files $WORKING_DIRECTORY/rulex.db
    remove_files $WORKING_DIRECTORY/upload/
    remove_files $WORKING_DIRECTORY/license.key
    remove_files $WORKING_DIRECTORY/license.lic
    rm -f "$WORKING_DIRECTORY/*.txt"
    rm -f "$WORKING_DIRECTORY/*.txt.gz"
    systemctl daemon-reload
    systemctl reset-failed
    echo "[√] Rulex has been uninstalled."
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