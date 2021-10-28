#! /bin/bash
LARTEST_VERSION=`curl -s https://raw.githubusercontent.com/wwhai/rulex/master/VERSION`
echo "LARTEST_VERSION is: ${LARTEST_VERSION}"
LOCAL_VERSION=$(cat ./VERSION)
echo "LOCAL_VERSION is: ${LARTEST_VERSION}"
if [ "${LARTEST_VERSION}X" == "${LOCAL_VERSION}X" ];then
    echo "Current is newest version!"
    exit 1
else
    read -r -p "Newest version found:[${LARTEST_VERSION}], Are You Sure upgrade? [Y/n] " input
    case $input in
        [yY][eE][sS]|[yY])
            echo "-- Download newest package"
            echo "-- Upgrade to [${LARTEST_VERSION}] finished"
            exit 1
            ;;

        [nN][oO]|[nN])
            exit 1
            ;;
        *)
            echo "Only support: [yY][eE][sS]|[yY] | [nN][oO]|[nN]"
            exit 1
            ;;
    esac
fi