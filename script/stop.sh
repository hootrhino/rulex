#!/bin/bash

#--------------------------
# 停止: Rulex
#--------------------------
stop() {
    PID=$(ps -a | grep "$1" | awk '{print $1}')
    if [ -z "$PID" ]; then
        echo "$1 not exists"
        exit 1
    else
        RUNNING=$(ps -p $PID | awk 'NR==2{print $1}')
        if [ -z "$RUNNING" ]; then
            echo "$1 not exists"
            exit 1
        else
            read -r -p "* $1 is running, pid is: $PID, Are your sure to kill it? [Y/y/N/n] " input
            case $input in
            [yY][eE][sS] | [yY])
                echo "* Warning: will excute 'kill -9 $PID'"
                kill -9 $PID
                echo "* $1 has been killed"
                exit 1
                ;;
            [nN][oO] | [nN])
                exit 1
                ;;
            *)
                exit 1
                ;;
            esac
        fi

    fi

}
#
# main
#
stop rulex
