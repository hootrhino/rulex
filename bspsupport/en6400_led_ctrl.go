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

package archsupport

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	__EN6400_GPIO231 string = "231" // GPIO.231
)
const (
	__EN6400_GPIO231_PATH = "/sys/class/gpio/gpio231/value"
)

// echo 231 > /sys/class/gpio/export
// echo out > /sys/class/gpio/gpio231/direction
// 熄灭LED:
//     echo 1 >/sys/class/gpio/gpio231/value
// 点亮LED:
//     echo 0 >/sys/class/gpio/gpio231/value

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "EN6400" {
		_EN6400_GPIOAllInit()
	}
}
func _EN6400_GPIOAllInit() {
	_EN6400_LedInit(__EN6400_GPIO231, "out")
}

func _EN6400_LedInit(Pin string, direction string) {
	//gpio export
	cmd := fmt.Sprintf("echo %s > /sys/class/gpio/export", Pin)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[EN6400B_GPIOInit] error", err)
		return
	}
	// gpio set direction
	cmd = fmt.Sprintf("echo %s > /sys/class/gpio/gpio%s/direction", direction, Pin)
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[EN6400B_GPIOInit] error", err)
	}
}

func EN6400_GPIO231Set(value int) error {
	return _EN6400_GPIOSet(__EN6400_GPIO231_PATH, value)
}

func EN6400_GPIO231Get() (int, error) {
	return _EN6400_GPIOGet(__EN6400_GPIO231_PATH)

}

func _EN6400_GPIOSet(gpioPath string, value int) error {
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
func _EN6400_GPIOGet(gpioPath string) (int, error) {
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
