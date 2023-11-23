#!/bin/bash

### BEGIN INIT INFO
# Provides:          rulex
# Required-Start:    $network $local_fs $remote_fs
# Required-Stop:     $network $local_fs $remote_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Rulex Service
# Description:       Rulex Service
### END INIT INFO

EXECUTABLE_PATH="/usr/local/rulex"
CONFIG_PATH="/usr/local/rulex.ini"
SERVICE_NAME="rulex"
WORKING_DIRECTORY="/usr/local/"
WAIT_TIME_SECONDS=3
CHECK_INTERVAL_SECONDS=1
PID_FILE="/var/run/$SERVICE_NAME.pid"
SCRIPT_PATH="/etc/init.d/rulex.sh"
PID_FILE="/var/run/rulex.pid"

log() {
    local level=$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}


install(){
    local source_dir="$PWD"
    local service_file="/etc/init.d/rulex.service"
    local executable="/usr/local/rulex"
    local working_directory="/usr/local/"
    local config_file="/usr/local/rulex.ini"
    local db_file="/usr/local/rulex.db"
cat > "$service_file" << EOL
#!/bin/sh $service_file
# Create Time: $(date +'%Y-%m-%d %H:%M:%S')
log() {
    local level=\$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}

start() {
    log INFO "Starting rulex..."
    nohup $executable run -config=$config_file nohup.log 2>&1 &
    echo $! > "$PID_FILE"
    log INFO "Starting rulex Finished"
}

stop() {
    # Check if rulex process is running
    if pgrep -x "rulex" > /dev/null; then
        pid=$(pgrep -x "rulex")
        log INFO "Killing rulex process with PID $pid"
        kill "$pid"
    else
        log INFO "rulex process is not running."
    fi
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

    mkdir -p $working_directory
    chmod +x $source_dir/rulex
    cp -rfp "$source_dir/rulex" "$executable"
    cp -rfp "$source_dir/rulex.ini" "$config_file"
    cp -rfp "$source_dir/license.lic" "$working_directory"
    cp -rfp "$source_dir/license.key" "$working_directory"
    chmod 777 $service_file
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

uninstall(){
    local working_directory="/usr/local"
    __remove_files /etc/systemd/system/rulex.service
    __remove_files $working_directory/rulex
    __remove_files $working_directory/rulex.ini
    __remove_files $working_directory/rulex.db
    __remove_files $working_directory/license.lic
    __remove_files $working_directory/license.key
    __remove_files $working_directory/upload/
    __remove_files $working_directory/*.txt
    __remove_files $working_directory/*.txt.gz
    log INFO "Rulex has been uninstalled."
}

start() {
    log INFO "Starting $SERVICE_NAME..."
    nohup $EXECUTABLE_PATH run -config $CONFIG_PATH >output.log 2>&1 &
    echo "$!" > "$PID_FILE" && log INFO "Service started."
}

stop() {
    log INFO "Stopping $SERVICE_NAME..."
    rm -f "$PID_FILE" && log INFO "PID file removed."
    pid=$(pgrep -x "$SERVICE_NAME")
    if [ -n "$pid" ]; then
        kill -15 "$pid" && log INFO "Process $pid (rulex) terminated."
        pkill -f "$SERVICE_NAME"
    else
        log INFO "Process rulex not found."
    fi
}

restart(){
    stop
    start
}

status() {
    log INFO "Checking $SERVICE_NAME status..."
    pid=$(pgrep -x "$SERVICE_NAME")
    if [ -n "$pid" ]; then
        log INFO "$SERVICE_NAME is running with Pid:${pid}"
    else
        log INFO "$SERVICE_NAME is not running."
    fi
}

case "$1" in
    install)
        install
    ;;
    start)
        start
    ;;
    restart)
        stop
        start
    ;;
    stop)
        stop
    ;;
    uninstall)
        uninstall
    ;;
    status)
        status
    ;;
    *)
        log ERROR "Usage: $0 {install|start|restart|stop|uninstall|status}"
        exit 1
    ;;
esac

exit 0
