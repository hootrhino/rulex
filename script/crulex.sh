#!/bin/bash
set +e
INSTALL_PATH="/usr/local/rulex"
#
#
#
HELP_TEXT="Please use './crulex [release-name] [version]' command. Such as './crulex rulex-x64linux V0.0.9'."
RELEASE=$1
VERSION=$2
if [ -z "$1" ]; then
    echo ">> Release arg missed. ${HELP_TEXT}"
    exit 1
fi
if [ "$1" == "releases" ]; then
    curl -s -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/i4de/rulex/releases |
        jq '[ .[].assets | .[] | .name]'
    exit 1
fi
if [ -z "$2" ]; then
    echo ">> Version arg missed. ${HELP_TEXT}"
    exit 1
fi
#
#
#
init_env() {
    if [ ! -d ${INSTALL_PATH} ]; then
        mkdir -p ${INSTALL_PATH}
    fi
}

#
#
#
add_service() {
    cat >/usr/lib/systemd/system/rulex.service <<-EOF
[Unit]
Description=RULEX Engine ${RELEASE}-${VERSION}
[Service]
User=root
WorkingDirectory=/usr/local/rulex/
TimeoutStartSec=5
ExecStart=/usr/local/rulex/rulex run -config=/usr/local/rulex/rulex.ini
ExecStop=echo "RULEX stop."
[Install]
WantedBy=multi-user.target
EOF
    echo "----------------------------------"
    echo ">> RULEX installed successfully!!! "
    echo ">> systemctl command list: "
    echo "* sudo systemctl start rulex"
    echo "* sudo systemctl enable rulex"
    echo "* sudo systemctl status rulex"
    echo "----------------------------------"
}
#
# install
#
install() {
    echo ">> Installing ${RELEASE}-${VERSION}"
    init_env
    create_temp_path
    cd ./_temp/
    URL=https://github.com/i4de/rulex/releases/download/${VERSION}/${RELEASE}-${VERSION}.zip
    echo ">> Download ${RELEASE}-${VERSION} from: ${URL}"
    wget -q --show-progress ${URL}
    unzip -o ${RELEASE}-${VERSION}.zip
    cp ./${RELEASE} ${INSTALL_PATH}
    cp ./conf/rulex.ini ${INSTALL_PATH}/rulex.ini
    add_service
    cd ../
    rm -rf ./_temp/
}
#
# create_temp_path
#
create_temp_path() {
    if [ ! -d "./_temp/" ]; then
        mkdir -p ./_temp/
    else
        rm -rf ./_temp/
        mkdir -p ./_temp/
    fi

}
#
# Install
#
install
