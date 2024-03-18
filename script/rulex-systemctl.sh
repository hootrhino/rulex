#!/bin/bash

SERVICE_NAME="rulex"
WORKING_DIRECTORY="/usr/local"
EXECUTABLE_PATH="$WORKING_DIRECTORY/$SERVICE_NAME"
CONFIG_PATH="$WORKING_DIRECTORY/$SERVICE_NAME.ini"

SERVICE_FILE="/etc/systemd/system/rulex.service"

STOP_SIGNAL="/var/run/rulex-stop.sinal"
UPGRADE_SIGNAL="/var/run/rulex-upgrade.lock"

SOURCE_DIR="$PWD"

install(){
cat > "$SERVICE_FILE" << EOL
[Unit]
Description=Rulex Daemon
After=network.target

[Service]
Environment="ARCHSUPPORT=EEKITH3"
WorkingDirectory=$WORKING_DIRECTORY
ExecStart=$EXECUTABLE_PATH run
ConditionPathExists=!/var/run/rulex-upgrade.lock
Restart=always
User=root
Group=root
StartLimitInterval=0
RestartSec=5
[Install]
WantedBy=multi-user.target
EOL
    chmod +x $SOURCE_DIR/rulex
    echo "[.] Copy $SOURCE_DIR/rulex to $WORKING_DIRECTORY."
    cp "$SOURCE_DIR/rulex" "$EXECUTABLE_PATH"
    echo "[.] Copy $SOURCE_DIR/rulex.ini to $WORKING_DIRECTORY."
    cp "$SOURCE_DIR/rulex.ini" "$config_file"
    echo "[.] Copy $SOURCE_DIR/license.key to /usr/local/license.key."
    cp "$SOURCE_DIR/license.key" "/usr/local/license.key"
    echo "[.] Copy $SOURCE_DIR/license.lic to /usr/local/license.lic."
    cp "$SOURCE_DIR/license.lic" "/usr/local/license.lic"
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
    remove_files "$SERVICE_FILE"
    remove_files "$WORKING_DIRECTORY/rulex"
    remove_files "$WORKING_DIRECTORY/rulex.ini"
    remove_files "$WORKING_DIRECTORY/rulex.db"
    remove_files "$WORKING_DIRECTORY/license.lic"
    remove_files "$WORKING_DIRECTORY/license.key"
    remove_files "$WORKING_DIRECTORY/rulex_internal_datacenter.db"
    remove_files "$WORKING_DIRECTORY/upload/"
    remove_files "$WORKING_DIRECTORY/rulexlog.txt"
    remove_files "$WORKING_DIRECTORY/rulex-daemon-log.txt"
    remove_files "$WORKING_DIRECTORY/rulex-recover-log.txt"
    remove_files "$WORKING_DIRECTORY/rulex-upgrade-log.txt"
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