#!/bin/sh /etc/rc.common
# rulex_daemon - Rulex daemon script for Linux

START=99
USE_PROCD=1

# Check if the service is disabled
[ -e /etc/config/rulex_daemon ] && . /etc/config/rulex_daemon

working_directory="./"

# Function to copy files to /usr/local
install_files() {
    cp "$working_directory"/rulex /usr/local/
    cp "$working_directory"/rulex.ini /usr/local/
}

# Function to uninstall the service
uninstall_files() {
    rm -f /usr/local/rulex
    rm -f /usr/local/rulex.ini
    rm -f /usr/local/rulex.db
    rm -f /usr/local/*.txt
    rm -f /usr/local/license.*
    rm -rf /usr/local/upload/
    rm -f /usr/local/*.txt.gz
}

start_service() {
    if [ "$DISABLED" -eq 0 ]; then
        procd_open_instance
        procd_set_param command /usr/local/rulex run -config /usr/local/rulex.ini
        procd_set_param respawn
        procd_set_param timeout 5  # 5 seconds timeout
        procd_close_instance
    else
        echo "Service is disabled. To enable, run: /etc/init.d/rulex_daemon enable"
    fi
}

stop_service() {
    procd_close_instance
}

# Function to disable the service
disable_service() {
    [ -e /etc/config/rulex_daemon ] && echo 'DISABLED=1' > /etc/config/rulex_daemon
    /etc/init.d/rulex_daemon stop
    /etc/init.d/rulex_daemon disable
}

# Function to enable the service
enable_service() {
    [ -e /etc/config/rulex_daemon ] && rm /etc/config/rulex_daemon
    /etc/init.d/rulex_daemon enable
}

# Function to uninstall the service
uninstall_service() {
    uninstall_files
    /etc/init.d/rulex_daemon stop
    /etc/init.d/rulex_daemon disable
    rm /etc/init.d/rulex_daemon
    rm /etc/config/rulex_daemon
    echo "Rulex uninstallation complete."
}

# Function to check service status
status_service() {
    if procd_status rulex_daemon > /dev/null; then
        echo "Rulex is running."
    else
        echo "Rulex is not running."
    fi
}

service_triggers() {
    procd_add_reload_trigger "rulex"
}

reload_service() {
    procd_send_signal rulex HUP
}

shutdown_service() {
    procd_close_instance
}

service_error() {
    procd_close_instance
    echo "Error starting rulex_daemon" >&2
}

run() {
    case "$1" in
        install)
            install_files
        ;;
        start)
            start_service
        ;;
        restart)
            reload_service
        ;;
        stop)
            stop_service
        ;;
        uninstall)
            uninstall_service
        ;;
        status)
            status_service
        ;;
        *)
            echo "Usage: $0 {install|start|restart|stop|uninstall|status}"
            exit 1
        ;;
    esac
}
