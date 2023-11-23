#!/bin/bash
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


log() {
    local level=$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}

install() {
    local source_dir="$PWD"
    local service_file="/etc/init.d/rulex.service"
    local executable="/usr/local/rulex"
    local working_directory="/usr/local/"
    local config_file="/usr/local/rulex.ini"
    local db_file="/usr/local/rulex.db"

    cat > "$service_file" << EOL
#!/bin/sh $service_file

START=180

USE_PROCD=1

start() {
    procd_open_instance
    procd_set_param command "$executable run -config=$config_file -db=$db_file"
    procd_set_param respawn 0
    procd_set_param stdout 1
    procd_set_param stderr 1
    procd_close_instance
}

stop() {
    service_stop "$executable"
}

restart() {
    stop
    start
}

status() {
    log INFO "Checking rulex status..."
    pid=$(pgrep -x "rulex")
    if [ -n "$pid" ]; then
        log INFO "rulex is running with Pid:${pid}"
    else
        log INFO "rulex is not running."
    fi
}

EOL

    mkdir -p "$working_directory"
    chmod +x "$source_dir/rulex"
    cp -rfp "$source_dir/rulex" "$executable"
    cp -rfp "$source_dir/rulex.ini" "$config_file"
    cp -rfp "$source_dir/license.lic" "$working_directory"
    cp -rfp "$source_dir/license.key" "$working_directory"

    chmod 777 "$service_file"
    "$service_file" enable

    if [ $? -eq 0 ]; then
        log INFO "Rulex service has been created and extracted."
    else
        log ERROR "Failed to create the Rulex service or extract files."
    fi
    exit 0
}

__remove_files() {
    local file=$1
    log INFO "Removing $file..."
    if [ -e "$file" ]; then
        if [ -d "$file" ]; then
            rm -rf "$file"
        else
            rm "$file"
        fi
        log INFO "$file removed."
    else
        log INFO "$file not found. No need to remove."
    fi
}

uninstall() {
    local service_file="$service_file.service"
    "$service_file" stop
    "$service_file" disable
    local working_directory="/usr/local"
    __remove_files /etc/systemd/system/rulex.service
    __remove_files "$working_directory/rulex" "$working_directory/rulex.ini" "$working_directory/rulex.db"
    __remove_files "$working_directory/license.lic" "$working_directory/license.key"
    __remove_files "$working_directory/upload/" "$working_directory/"*.txt "$working_directory/"*.txt.gz
    log INFO "Rulex has been uninstalled."
}

start() {
    $service_file start
}

restart() {
    $service_file restart
}

stop() {
    $service_file stop
}

status() {
    $service_file running
}

main() {
    case "$1" in
        install | start | restart | stop | uninstall | create_user | status | openwrt)
            $1
        ;;
        *)
            log ERROR "Invalid command: $1"
            echo "[?] Usage: $0 <install|start|restart|stop|uninstall|status>"
            exit 1
        ;;
    esac
    exit 0
}

if [ "$(id -u)" != "0" ]; then
    log ERROR "This script must be run as root"
    exit 1
else
    main "$1"
fi
