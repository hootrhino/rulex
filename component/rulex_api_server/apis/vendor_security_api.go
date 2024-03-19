// Copyright (C) 2023 wwhai
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

package apis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 获取一机一密
*
 */
type LocalLicenseVo struct {
	DeviceID          string `json:"device_id"`          // 设备生产序列号
	AuthorizeAdmin    string `json:"authorize_admin"`    // 证书签发人
	AuthorizePassword string `json:"authorize_password"` // 证书签发人密钥
	BeginAuthorize    int64  `json:"begin_authorize"`    // 证书授权开始时间
	EndAuthorize      int64  `json:"end_authorize"`      // 证书授权结束时间
	MAC               string `json:"mac"`                // 设备硬件MAC地址，一般取以太网卡
	License           string `json:"license"`            // 公钥, 发给用户设备
}

// 00001 & rhino & hoot & FF:FF:FF:FF:FF:FF & 0 & 0
func ParseAuthInfo(info string) (LocalLicenseVo, error) {
	LocalLicense := LocalLicenseVo{}
	ss := strings.Split(info, "&")
	if len(ss) == 6 {
		BeginAuthorize, err1 := strconv.ParseInt(ss[4], 10, 64)
		if err1 != nil {
			return LocalLicense, err1
		}
		EndAuthorize, err2 := strconv.ParseInt(ss[5], 10, 64)
		if err2 != nil {
			return LocalLicense, err2
		}
		LocalLicense.DeviceID = ss[0]
		LocalLicense.AuthorizeAdmin = ss[1]
		LocalLicense.AuthorizePassword = ss[2]
		LocalLicense.MAC = ss[3]
		LocalLicense.BeginAuthorize = BeginAuthorize
		LocalLicense.EndAuthorize = EndAuthorize
		return LocalLicense, nil
	}
	return LocalLicense, fmt.Errorf("failed parse:%s", info)
}

/*
*
* 获取证书
*
 */
func GetVendorKey(c *gin.Context, ruleEngine typex.RuleX) {
	LocalLicenseVo := LocalLicenseVo{
		DeviceID:          "00001",
		AuthorizeAdmin:    "rhino",
		AuthorizePassword: "hoot",
		MAC:               "FF:FF:FF:FF:FF:FF",
		BeginAuthorize:    0,
		EndAuthorize:      0,
		License:           "DefaultLicense",
	}
	c.JSON(common.HTTP_OK, common.OkWithData(LocalLicenseVo))
}
