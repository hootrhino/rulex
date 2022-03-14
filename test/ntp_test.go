package test

import (
	"testing"

	"github.com/wwhai/ntp"
)

func Test_china_ntp(t *testing.T) {
	time, err := ntp.Time("0.cn.pool.ntp.org")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
}
