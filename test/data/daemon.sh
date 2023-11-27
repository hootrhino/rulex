#!/bin/sh
# Create Time: 2023-11-27 14:59:06

WORKING_DIRECTORY="/usr/local"
PID_FILE="/var/run/rulex.pid"
EXECUTABLE_PATH="$WORKING_DIRECTORY/rulex"
CONFIG_PATH="$WORKING_DIRECTORY/rulex.ini"

log() {
    local level=$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}

start() {
    rm -f /var/run/rulex-stop.sinal
    pid=$(pgrep -x -n -f "/usr/local/rulex run -config=/usr/local/rulex.ini")
    if [ -n "$pid" ]; then
        log INFO "rulex is running with Pid:${pid}"
        exit 0
    fi
    daemon &
    exit 0
}

stop() {
    echo "1" > /var/run/rulex-stop.sinal
    if pgrep -x "rulex" > /dev/null; then
        log INFO "rulex process is running. Killing it..."
        pkill -x "rulex"
        log INFO "rulex process has been killed."
    else
        log WARNING "rulex process is not running."
    fi
}

restart() {
    stop
    sleep 1
    start
}

status() {
    log INFO "Checking rulex status."
    pid=$(pgrep -x -n "rulex")
    if [ -n "$pid" ]; then
        log INFO "rulex is running with Pid:${pid}"
    else
        log INFO "rulex is not running."
    fi
}

daemon() {
    while true; do
        if pgrep -x "rulex" > /dev/null; then
            sleep 3
            continue
        fi
        if ! pgrep -x "rulex" > /dev/null; then
            if [ -e "/var/run/rulex-upgrade.lock" ]; then
                log INFO "File /var/run/rulex-upgrade.lock exists. May upgrade now."
                sleep 2
                continue
            elif [ -e "/var/run/rulex-stop.sinal" ]; then
                log INFO "/var/run/rulex-stop.sinal file found. Exiting."
                exit 0
            else
                log WARNING "Detected that rulex process is interrupted. Restarting..."
                /usr/local/rulex run -config=/usr/local/rulex.ini
                log WARNING "Detected that rulex process has Restarted."
            fi
        fi
        sleep 4
    done
}

case "$1" in
    start)
        start
    ;;
    restart)
        restart
    ;;
    stop)
        stop
    ;;
    status)
        status
    ;;
    *)
        log ERROR "Usage: $0 {start|restart|stop|status}"
        exit 1
    ;;
esac

