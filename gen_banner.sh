#! /bin/bash
VERSION="$(git describe --tags $(git rev-list --tags --max-count=1))"
HASH=`git rev-list --tags --max-count=1`
cat >./VERSION <<EOF
$VERSION
EOF
cat >./conf/banner.txt <<EOF
 -------------------------------------------------------------
|                ____ _  _ _    ____ _  _                      |
|                |__/ |  | |    |___  \/           ------      |
|                |  \ |__| |___ |___ _/\_       ------         |
|                                            ------            |
|                                         ------               |
|* Version: ${VERSION}-${HASH:0:15}                             |
|* Document: https://wwhai.github.io/rulex_doc_html/index.html |
|                                                              |
 --------------------------------------------------------------
EOF
echo "Generate Banner Susseccfully"
