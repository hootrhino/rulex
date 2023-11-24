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

SERVICE_NAME="rulex"
WORKING_DIRECTORY="/usr/local"
EXECUTABLE_PATH="$WORKING_DIRECTORY/$SERVICE_NAME"
CONFIG_PATH="$WORKING_DIRECTORY/$SERVICE_NAME.ini"

PID_FILE="/var/run/$SERVICE_NAME.pid"
PID_FILE="/var/run/rulex.pid"
SERVICE_FILE="/etc/init.d/rulex.service"
log() {
    local level=$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}


install(){
    local source_dir="$PWD"
    local db_file="/usr/local/rulex.db"
cat > "$SERVICE_FILE" << EOL
#!/bin/sh
# Create Time: $(date +'%Y-%m-%d %H:%M:%S')

WORKING_DIRECTORY="/usr/local"
PID_FILE="/var/run/rulex.pid"
EXECUTABLE_PATH="\$WORKING_DIRECTORY/rulex"
CONFIG_PATH="\$WORKING_DIRECTORY/rulex.ini"

log() {
    local level=\$1
    shift
    echo "[\$level] \$(date +'%Y-%m-%d %H:%M:%S') - \$@"
}

start() {
    pid=\$(pgrep -x -n "rulex")
    if [ -n "\$pid" ]; then
        log INFO "rulex is running with Pid:\${pid}"
        exit 0
    fi
    log INFO "Starting rulex."
    cd \$WORKING_DIRECTORY
    echo \$! > "\$PID_FILE"
    daemon > rulex-daemon-log.txt 2>&1 &
    log INFO "Starting rulex Finished"
}

stop() {
    if pgrep -x -n "rulex" > /dev/null; then
        pid=\$(pgrep -x -n "rulex")
        kill "\$pid"
        log INFO "Killing rulex process with PID \$pid"
    else
        log INFO "rulex process is not running."
    fi
}

restart() {
    stop
    start
}

status() {
    log INFO "Checking rulex status."
    pid=\$(pgrep -x -n "rulex")
    if [ -n "\$pid" ]; then
        log INFO "rulex is running with Pid:\${pid}"
    else
        log INFO "rulex is not running."
    fi
}

daemon(){
    while true; do
        $EXECUTABLE_PATH run -config=\$CONFIG_PATH
        wait \$!
        sleep 3
        if [ ! -f "\$PID_FILE" ]; then
            log INFO "Pid File Remove detected, RULEX Daemon Exit."
            exit 0
        fi
    done
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
    log INFO "Copy rulex to $WORKING_DIRECTORY"
    cp -rfp "$source_dir/rulex" "$EXECUTABLE_PATH"

    log INFO "Copy rulex.ini to $WORKING_DIRECTORY"
    cp -rfp "$source_dir/rulex.ini" "$CONFIG_PATH"

    log INFO "Copy license.lic to $WORKING_DIRECTORY"
    cp -rfp "$source_dir/license.lic" "$WORKING_DIRECTORY/"

    log INFO "Copy license.key to $WORKING_DIRECTORY"
    cp -rfp "$source_dir/license.key" "$WORKING_DIRECTORY/"
    __add_to_rc_local
    chmod 777 $SERVICE_FILE
    if [ $? -eq 0 ]; then
        log INFO "Rulex service has been created and extracted."
    else
        log ERROR "Failed to create the Rulex service or extract files."
    fi
    exit 0
}

__remove_files() {
    local file=$1
    log INFO "Removing $file."
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
__remove_from_rc_local() {
    local rc_local_path="/etc/rc.local"
    if [ ! -f "$rc_local_path" ]; then
        log ERROR "Error: /etc/rc.local does not exist. Check your system configuration."
        return 1
    fi
    if ! grep -qF "$SERVICE_FILE start" "$rc_local_path"; then
        log INFO "Script not found in /etc/rc.local. No changes made."
        return 0
    fi
    sed -i "\|$SERVICE_FILE start|d" "$rc_local_path"
    log INFO "Script removed from /etc/rc.local."
    return 0
}

__add_to_rc_local() {
    local rc_local_path="/etc/rc.local"
    if [ ! -f "$rc_local_path" ]; then
        log INFO "Error: /etc/rc.local does not exist. Create the file manually or check your system configuration."
        return 1
    fi
    if grep -qF "$SERVICE_FILE start" "$rc_local_path"; then
        log INFO "Script already present in /etc/rc.local. No changes made."
        return 0
    fi
    local last_line_number=$(awk '/^[^#[:space:]]/{n=$0} END{print NR}' "$rc_local_path")
    if [ -n "$last_line_number" ]; then
        sed -i "${last_line_number}i $SERVICE_FILE start" "$rc_local_path"
    else
        echo "$SERVICE_FILE start" >> "$rc_local_path"
    fi
    log INFO "Script added to /etc/rc.local."
    return 0
}

uninstall(){
    if [ -e "$SERVICE_FILE" ]; then
        $SERVICE_FILE stop
    fi
    __remove_files "$PID_FILE"
    __remove_files "$SERVICE_FILE"
    __remove_files "$WORKING_DIRECTORY/rulex"
    __remove_files "$WORKING_DIRECTORY/rulex.ini"
    __remove_files "$WORKING_DIRECTORY/rulex.db"
    __remove_files "$WORKING_DIRECTORY/license.lic"
    __remove_files "$WORKING_DIRECTORY/license.key"
    __remove_files "$WORKING_DIRECTORY/RULEX_INTERNAL_DATACENTER.db"
    __remove_files "$WORKING_DIRECTORY/upload/"
    __remove_files "$WORKING_DIRECTORY/rulex-nohup-log.txt"
    __remove_files "$WORKING_DIRECTORY/rulexlog.txt"
    __remove_files "$WORKING_DIRECTORY/rulex-recover-log.txt"
    __remove_from_rc_local
    log INFO "Rulex has been uninstalled."
}

start() {
    $SERVICE_FILE start
}

restart() {
    $SERVICE_FILE restart
}

stop() {
    $SERVICE_FILE stop
}

status() {
    $SERVICE_FILE status
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
