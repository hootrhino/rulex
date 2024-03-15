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

package globalregistry

import "time"

type TimerWheel struct {
	interval time.Duration // 每个槽位的时间间隔
	ticker   *time.Ticker
	slots    []chan func() // 时间槽
	size     int           // 时间槽数量
	current  int           // 当前指针位置
	quit     chan struct{} // 退出信号
}

func NewTimerWheel(interval time.Duration, size int) *TimerWheel {
	tw := &TimerWheel{
		interval: interval,
		ticker:   time.NewTicker(interval),
		size:     size,
		slots:    make([]chan func(), size),
		quit:     make(chan struct{}),
	}
	for i := range tw.slots {
		tw.slots[i] = make(chan func(), 100) // 每个槽位最多可存储 100 个任务
	}
	// go tw.run()
	return tw
}

func (tw *TimerWheel) AddTimer(d time.Duration, f func()) {
	index := (tw.current + int(d/tw.interval)) % tw.size
	tw.slots[index] <- f
}

func (tw *TimerWheel) Run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.onTicker()
		case <-tw.quit:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimerWheel) onTicker() {
	index := tw.current
	tw.current = (tw.current + 1) % tw.size
	for {
		select {
		case f := <-tw.slots[index]:
			f()
		default:
			return
		}
	}
}

func (tw *TimerWheel) Stop() {
	close(tw.quit)
}

// func main() {
// 	tw := NewTimerWheel(time.Second, 10)

// 	for i := 0; i < 20; i++ {
// 		delay := time.Duration(i) * time.Second
// 		tw.AddTimer(delay, func() {
// 			fmt.Println("Task executed after", delay)
// 		})
// 	}

// 	time.Sleep(22 * time.Second)
// 	tw.Stop()
// }
