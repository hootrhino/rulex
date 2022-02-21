#!/bin/bash
INSTALL_PATH="/usr/local/rulex"
init_env() {
    if [ ! -d ${INSTALL_PATH} ]; then
        mkdir -p /usr/local/rulex
    else
        echo "Path '/usr/local/rulex' exists."
    fi
}
echo "Installing rulex...."
init_env
cp ../rulex/rulex /usr/local/rulex
cp ../rulex/conf/rulex.ini /usr/local/rulex/rulex.ini
cp ../rulex/script/rulex.service /usr/lib/systemd/system/rulex.service