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
WORKING_DIRECTORY="/usr/local"
WAIT_TIME_SECONDS=3
CHECK_INTERVAL_SECONDS=1
PID_FILE="/var/run/$SERVICE_NAME.pid"
SCRIPT_PATH="/etc/init.d/rulex.sh"
PID_FILE="/var/run/rulex.pid"
service_file="/etc/init.d/rulex.service"
log() {
    local level=$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}


install(){
    local source_dir="$PWD"
    local executable="/usr/local/rulex"
    local config_file="/usr/local/rulex.ini"
    local db_file="/usr/local/rulex.db"
cat > "$service_file" << EOL
#!/bin/sh
# Create Time: $(date +'%Y-%m-%d %H:%M:%S')

WORKING_DIRECTORY="/usr/local"
PID_FILE="/var/run/rulex.pid"
executable="/usr/local/rulex"
config_file="/usr/local/rulex.ini"

log() {
    local level=\$1
    shift
    echo "[\$level] \$(date +'%Y-%m-%d %H:%M:%S') - \$@"
}

start() {
    log INFO "Starting rulex..."
    cd \$WORKING_DIRECTORY
    nohup \$executable run -config=\$config_file > run-nohup-log.txt 2>&1 &
    echo \$! > "\$PID_FILE"
    log INFO "Starting rulex Finished"
}

stop() {
    # Check if rulex process is running
    if pgrep -x "rulex" > /dev/null; then
        pid=\$(pgrep -x "rulex")
        log INFO "Killing rulex process with PID \$pid"
        kill "\$pid"
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
    pid=\$(pgrep -x "rulex")
    if [ -n "\$pid" ]; then
        log INFO "rulex is running with Pid:\${pid}"
    else
        log INFO "rulex is not running."
    fi
}

case "\$1" in
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
    status)
        status
    ;;
    *)
        log ERROR "Usage: \$0 {start|restart|stop|status}"
        exit 1
    ;;
esac

EOL

    mkdir -p $WORKING_DIRECTORY
    chmod +x $source_dir/rulex
    cp -rfp "$source_dir/rulex" "$executable"
    cp -rfp "$source_dir/rulex.ini" "$config_file"
    cp -rfp "$source_dir/license.lic" "$WORKING_DIRECTORY/"
    cp -rfp "$source_dir/license.key" "$WORKING_DIRECTORY/"
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
    if [ -e "$service_file" ]; then
        $service_file stop
        $service_file disable
    fi
    __remove_files $service_file
    __remove_files "$WORKING_DIRECTORY/rulex"
    __remove_files "$WORKING_DIRECTORY/rulex.ini"
    __remove_files "$WORKING_DIRECTORY/rulex.db"
    __remove_files "$WORKING_DIRECTORY/license.lic"
    __remove_files "$WORKING_DIRECTORY/license.key"
    __remove_files "$WORKING_DIRECTORY/RULEX_INTERNAL_DATACENTER.db"
    __remove_files "$WORKING_DIRECTORY/LICENSE"
    __remove_files "$WORKING_DIRECTORY/md5.sum"
    __remove_files "$WORKING_DIRECTORY/upload/"
    __remove_files "$WORKING_DIRECTORY/*.txt"
    __remove_files "$WORKING_DIRECTORY/*.txt.gz"
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
    $service_file status
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
