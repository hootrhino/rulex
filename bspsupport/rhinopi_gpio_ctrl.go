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

package archsupport

import (
	"fmt"
	"os"
)

/*
*
* GPIO 表
*
 */
const (
	__h3_GPIO_PATH = "/sys/class/gpio/gpio%v/value"
	// __h3_DO1       = "/sys/class/gpio/gpio6/value"  // gpio6
	// __h3_DO2       = "/sys/class/gpio/gpio7/value"  // gpio7
	// __h3_DI1       = "/sys/class/gpio/gpio8/value"  // gpio8
	// __h3_DI2       = "/sys/class/gpio/gpio9/value"  // gpio9
	// __h3_DI3       = "/sys/class/gpio/gpio10/value" // gpio10
)

/*
*
* 新版本的文件读取形式获取GPIO状态
*
 */
func EEKIT_GPIOGetDO1() (int, error) {
	return EEKIT_GPIOGetByFile(6)
}
func EEKIT_GPIOGetDO2() (int, error) {
	return EEKIT_GPIOGetByFile(7)
}
func EEKIT_GPIOGetDI1() (int, error) {
	return EEKIT_GPIOGetByFile(8)
}
func EEKIT_GPIOGetDI2() (int, error) {
	return EEKIT_GPIOGetByFile(9)
}
func EEKIT_GPIOGetDI3() (int, error) {
	return EEKIT_GPIOGetByFile(10)
}
func EEKIT_GPIOGetByFile(pin byte) (int, error) {
	return __GPIOGet(fmt.Sprintf(__h3_GPIO_PATH, pin))
}

func __GPIOGet(gpioPath string) (int, error) {
	bites, err := os.ReadFile(gpioPath)
	if err != nil {
		return 0, err
	}
	if len(bites) > 0 {
		if bites[0] == '0' || bites[0] == 48 {
			return 0, nil
		}
		if bites[1] == '1' || bites[0] == 49 {
			return 1, nil
		}
	}
	return 0, fmt.Errorf("read gpio value failed: %s, value: %v", gpioPath, bites)
}

// Set

func EEKIT_GPIOSetDO1(value int) error {
	return EEKIT_GPIOSetByFile(6, value)
}
func EEKIT_GPIOSetDO2(value int) error {
	return EEKIT_GPIOSetByFile(7, value)
}
func EEKIT_GPIOSetDI1(value int) error {
	return EEKIT_GPIOSetByFile(8, value)
}
func EEKIT_GPIOSetDI2(value int) error {
	return EEKIT_GPIOSetByFile(9, value)
}
func EEKIT_GPIOSetDI3(value int) error {
	return EEKIT_GPIOSetByFile(10, value)
}

func EEKIT_GPIOSetByFile(pin, value int) error {
	return __GPIOSet(fmt.Sprintf(__h3_GPIO_PATH, pin), value)
}

func __GPIOSet(gpioPath string, value int) error {
	if value == 1 {
		err := os.WriteFile(gpioPath, []byte{'1'}, 0644)
		if err != nil {
			return err
		}
	}
	if value == 0 {
		err := os.WriteFile(gpioPath, []byte{'0'}, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
