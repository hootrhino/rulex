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

package hwportmanager

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

// 检查串口设备是否被系统进程占用
func CheckSerialBusy(serialDevice string) error {
	if runtime.GOOS == "linux" {
		_, err := syscall.Open(serialDevice, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0)
		if err != nil {
			pathErr, ok := err.(*os.PathError)
			if ok {
				if pathErr.Err.(syscall.Errno) == syscall.EBUSY {
					return fmt.Errorf("serial port %s is in use", serialDevice)
				}
			}
		}
	}
	return nil
}
