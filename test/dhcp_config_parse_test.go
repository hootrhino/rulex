// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test

import (
	"os"
	"testing"

	"github.com/hootrhino/go-dhcpd-leases"
)

// go test -timeout 30s -run ^Test_parse_linux_dhcp_leases github.com/hootrhino/rulex/test -v -count=1
func Test_parse_linux_dhcp_leases(t *testing.T) {
	f, err := os.Open("./data/dhcp.leases")
	if err != nil {
		t.Fatal(err)
	}
	leases := leases.Parse(f)
	for _, lease := range leases {
		t.Log(lease.IP.String(), lease.ClientHostname, lease.Hardware.MAC)
	}
	// t.Log(leases)
}
