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

package interclock

import (
	"time"
)

type InterClockTask func(any) error

/*
* 这是未来要实现的一个内部时间计数器。
* 时钟, 刻度：[1|2|3|4|5|6|7|8|9|10|11|12]
*
 */
type InterClock struct {
	Interval time.Duration    // 轮转圈数, 默认为10
	Ticker   *time.Ticker     // 时钟指针计时器
	Period   int              // 旋转一周的时间, 默认为 T= 1/10秒=100毫秒一圈
	Tasks    []InterClockTask // 任务列表
}
