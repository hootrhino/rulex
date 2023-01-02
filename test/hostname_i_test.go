package test

import (
	"testing"

	"github.com/i4de/rulex/utils"
)

// go test -timeout 30s -run ^Test_HostNameI github.com/i4de/rulex/test -v -count=1
func Test_HostNameI(t *testing.T) {
	// 172.30.211.225
	ip, err := utils.HostNameI()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}
