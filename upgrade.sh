#! /bin/bash
LARTEST_VERSION=`curl -s https://raw.githubusercontent.com/wwhai/rulex/master/VERSION`
echo "LARTEST_VERSION is: ${LARTEST_VERSION}"
LOCAL_VERSION=$(<./VERSION)
echo "LOCAL_VERSION is: ${LARTEST_VERSION}"
if [ "\"$LARTEST_VERSION"\" = "\"$LOCAL_VERSION"\" ];then
    echo "Current is newest version!"
    exit
else
    echo "Newest version found:${LARTEST_VERSION}"
fi