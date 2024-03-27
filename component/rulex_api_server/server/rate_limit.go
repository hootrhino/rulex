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
	"context"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var timeout int32 = 3000

func decreaseRateTime() {
	atomic.AddInt32(&timeout, -1)
}

func reInitRateTime() {
	atomic.StoreInt32(&timeout, 300)
}
func StartRateLimiter(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if atomic.LoadInt32(&timeout) > 0 {
				decreaseRateTime()
			}
		case <-ctx.Done():
			return // Stop the rate limiter when the context is cancelled
		}
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if atomic.LoadInt32(&timeout) <= 0 {
			reInitRateTime()
			c.Next()
		} else {
			c.AbortWithStatusJSON(429, map[string]interface{}{
				"code": 4001,
				"msg":  "Excessive operating frequency!",
			})
		}
	}
}
