#! /bin/bash
VERSION="$(git describe --tags $(git rev-list --tags --max-count=1))"
HASH=$(git rev-list --tags --max-count=1)

#######################################################################
## Gen Version
#######################################################################
cat >./typex/version.go <<EOF
//
// Warning:
//   This file is generated by go compiler, don't change it!!!
//   Build on: $(cat /etc/issue)
//
package typex

import "fmt"

type Version struct {
	Version     string
	ReleaseTime string
}

func (v Version) String() string {
	return fmt.Sprintf("{\"releaseTime\":\"%s\",\"version\":\"%s\"}", v.ReleaseTime, v.Version)
}

var DefaultVersion = Version{
	Version:   \`${VERSION}\`,
	ReleaseTime: "$(echo $(date "+%Y-%m-%d %H:%M:%S"))",
}
var Banner = \`
 **  Welcome to RULEX framework world <'_'>
**   Version: ${VERSION}-${HASH:0:15}
 **  Document: https://rulex.pages.dev
\`
EOF

echo "Generate Version Susseccfully"

