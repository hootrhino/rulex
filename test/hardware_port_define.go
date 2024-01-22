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

import "testing"

type BspPort struct {
	IFace  []string // 以太网卡
	WlFace []string // 无线网卡
	Audio  []string // 声卡
	Video  []string // 显卡
}

// go test -timeout 30s -run ^Test_gen_hd_define github.com/hootrhino/rulex/test -v -count=1
func Test_gen_hd_define(t *testing.T) {

}
