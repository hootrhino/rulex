package archsupport

/*
*
* Linux 特定实现
*
 */

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

//-----------------------------------------------
// 这是EEKIT网关的DI-DO支持库
//-----------------------------------------------
/*
    pins map

	DO1 -> PA6
	DO2 -> PA7
	DI1 -> PA8
	DI2	-> PA9
	DI3 -> PA10
*/
const (
	eekit_DO1 string = "6"
	eekit_DO2 string = "7"

	eekit_DI1 string = "8"
	eekit_DI2 string = "9"
	eekit_DI3 string = "10"
)

const (
	eekit_Out string = "out"
	eekit_In  string = "in"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "EEKITH3" {
		_EEKIT_GPIOAllInit()
	}
}

/*
explain:init all gpio
*/
func _EEKIT_GPIOAllInit() int {
	gpio6 := "/sys/class/gpio/gpio6/value"
	gpio7 := "/sys/class/gpio/gpio7/value"
	gpio8 := "/sys/class/gpio/gpio8/value"
	gpio9 := "/sys/class/gpio/gpio9/value"
	gpio10 := "/sys/class/gpio/gpio10/value"
	_, err1 := os.Stat(gpio6)
	_, err2 := os.Stat(gpio7)
	_, err3 := os.Stat(gpio8)
	_, err4 := os.Stat(gpio9)
	_, err5 := os.Stat(gpio10)
	if err1 != nil {
		_EEKIT_GPIOInit(eekit_DO1, eekit_Out)
	}
	if err2 != nil {
		_EEKIT_GPIOInit(eekit_DO2, eekit_Out)
	}
	if err3 != nil {
		_EEKIT_GPIOInit(eekit_DI1, eekit_In)
	}
	if err4 != nil {
		_EEKIT_GPIOInit(eekit_DI2, eekit_In)
	}
	if err5 != nil {
		_EEKIT_GPIOInit(eekit_DI3, eekit_In)
	}
	// 返回值无用
	return 1
}

/*
explain:init gpio
Pin: gpio pin
EnDir:gpio direction in or out
*/
func _EEKIT_GPIOInit(Pin string, EnDir string) {
	//gpio export
	cmd := fmt.Sprintf("echo %s > /sys/class/gpio/export", Pin)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[EEKIT_GPIOInit] error", err)
		return
	}
	//gpio set direction
	cmd = fmt.Sprintf("echo %s > /sys/class/gpio/gpio%s/direction", EnDir, Pin)
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[EEKIT_GPIOInit] error", err)
	}
}

/*
explain:set gpio
Pin: gpio pin
Value:gpio level 1 is high 0 is low
*/
func EEKIT_GPIOSet(pin, value int) (bool, error) {
	cmd := fmt.Sprintf("echo %d > /sys/class/gpio/gpio%d/value", value, pin)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[EEKIT_GPIOSet] error", err)
		return false, err
	}
	v, e := EEKIT_GPIOGet(pin)
	if e != nil {
		return false, e
	}
	return v == value, nil
}

/*
explain:read gpio
Pin: gpio pin
return:1 is high 0 is low
*/
func EEKIT_GPIOGet(pin int) (int, error) {
	cmd := fmt.Sprintf("cat /sys/class/gpio/gpio%d/value", pin)
	Value, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return -1, err
	}
	if len(Value) < 1 {
		return -1, errInvalidLen
	}
	if Value[0] == '0' {
		return 0, nil
	}
	if Value[0] == '1' {
		return 1, nil
	}
	return -1, errInvalidValue
}
