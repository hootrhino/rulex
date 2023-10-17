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

package server

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var timeout = 0
var __lock = sync.Mutex{}

func DecreaseRateTime() {
	__lock.Lock()
	defer __lock.Unlock()
	timeout--
}
func ReInitRateTime() {
	__lock.Lock()
	defer __lock.Unlock()
	timeout = 300
}
func StartRateLimiter() {
	for {
		if timeout == 0 {
			ReInitRateTime()
		}
		DecreaseRateTime()
		time.Sleep(100 * time.Millisecond)
	}
}
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout == 0 {
			c.Next()
			ReInitRateTime()
		} else {
			c.AbortWithStatusJSON(400, map[string]interface{}{
				"code": 4001,
				"msg":  "Excessive operating frequency!",
			})
		}
	}
}
