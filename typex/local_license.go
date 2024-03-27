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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package typex

import "time"

type LocalLicense struct {
	DeviceID          string `json:"device_id"`          // 设备生产序列号
	AuthorizeAdmin    string `json:"authorize_admin"`    // 证书签发人
	AuthorizePassword string `json:"authorize_password"` // 证书签发人密钥
	BeginAuthorize    int64  `json:"begin_authorize"`    // 证书授权开始时间
	EndAuthorize      int64  `json:"end_authorize"`      // 证书授权结束时间
	MAC               string `json:"mac"`                // 设备硬件MAC地址，一般取以太网卡
	License           string `json:"license"`            // 公钥, 发给用户设备
}

func (d LocalLicense) ValidateTime() bool {
	Now := time.Now().UnixNano()
	V := d.EndAuthorize - Now
	if (d.BeginAuthorize > Now) && (V <= 0) {
		return false
	}
	return true
}
