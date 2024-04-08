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
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/typex"
)

func CheckLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkLicense(c)
	}
}

func checkLicense(c *gin.Context) {
	if len(typex.License.License) == 0 {
		c.AbortWithStatusJSON(400, map[string]interface{}{
			"code": 4001,
			"msg":  "Invalid license!",
		})
		c.Abort()
		return
	}
	c.Next()
}
