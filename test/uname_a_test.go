package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-ini/ini"
	"github.com/hootrhino/rulex/glogger"
)

func ReadOSRelease(configfile string) *ini.File {
	cfg, err := ini.Load(configfile)
	if err != nil {
		glogger.GLogger.Fatal("Fail to read file: ", err)
	}
	return cfg
}

// NAME="Ubuntu"
// VERSION="20.04.2 LTS (Focal Fossa)"
// ID=ubuntu
// ID_LIKE=debian
// PRETTY_NAME="Ubuntu 20.04.2 LTS"
// VERSION_ID="20.04"
// HOME_URL="https://www.ubuntu.com/"
// SUPPORT_URL="https://help.ubuntu.com/"
// BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
// PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
// VERSION_CODENAME=focal
// UBUNTU_CODENAME=focal
func Test_uname_a(t *testing.T) {
	OSInfo := ReadOSRelease("/etc/os-release")
	fmt.Print(OSInfo.Section(""))
}
func Test_Mktemp_a(t *testing.T) {
	t.Log(os.MkdirTemp("./data", ""))
}
