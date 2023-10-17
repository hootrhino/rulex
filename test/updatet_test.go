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

package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// go test -timeout 30s -run ^TestApp_Hot_update github.com/hootrhino/rulex/test -v -count=1
func TestApp_Hot_update(t *testing.T) {
	fmt.Println(os.Args[0])
	s := gin.Default()
	s.GET("/rulex", downloadRulex)
	// s.StaticFile("/rulex", "./rulex")
	s.Run(":8088")
}

var (
	supportedArch = Strings{"amd64", "arm64"}
	supportedOS   = Strings{"linux", "windows"}
)

type Strings []string

func (ss Strings) Has(s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func downloadRulex(c *gin.Context) {
	var (
		os   = c.Query("os")
		arch = c.Query("arch")
		// version = c.Query("version")
	)

	// 检查参数
	if !supportedArch.Has(arch) {
		c.Status(http.StatusNotFound)
		return
	}
	if !supportedOS.Has(os) {
		c.Status(http.StatusNotFound)
		return
	}

	// TODO: 这里测试，暂时不检查version
	// 当请求字段没有version时，当作latest（最新版本）处理
	// 当version字段有效且大于等于最新版本时，不返回内容

	c.Writer.Header().Add("Rulex-Version", "v1.5.x")
	c.Writer.Header().Add("Rulex-MD5", "cc21c06617b95eabc812b82ffad2e9a8")
	c.File("./rulex_new")
}
