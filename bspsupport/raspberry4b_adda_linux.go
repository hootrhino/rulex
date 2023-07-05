package archsupport

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

/*
*
* 树莓派文档可以参考这里:https://www.raspberrypi.com/documentation/computers/raspberry-pi.html
*
 */

const (
	//-----------------------------------------------
	// 注意: 该定义使用的是树莓派GPIO的【物理编号】
	//-----------------------------------------------
	// 4个输出
	raspi_DO1 string = "11" // GPIO.0
	raspi_DO2 string = "12" // GPIO.1
	raspi_DO3 string = "13" // GPIO.2
	raspi_DO4 string = "15" // GPIO.3

	// 4个输入
	raspi_DI1 string = "16" // GPIO.4
	raspi_DI2 string = "18" // GPIO.5
	raspi_DI3 string = "22" // GPIO.6
	raspi_DI4 string = "29" // GPIO.21
)

const (
	raspi_Out string = "out"
	raspi_In  string = "in"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "RPI4B" {
		_RASPI4B_GPIOAllInit()
	}
}
func _RASPI4B_GPIOAllInit() {
	// 初始化输入
	_RASPI4B_GPIOInit(raspi_DO1, raspi_Out)
	_RASPI4B_GPIOInit(raspi_DO2, raspi_Out)
	_RASPI4B_GPIOInit(raspi_DO3, raspi_Out)
	_RASPI4B_GPIOInit(raspi_DO4, raspi_Out)
	// 初始化输出
	_RASPI4B_GPIOInit(raspi_DI1, raspi_In)
	_RASPI4B_GPIOInit(raspi_DI2, raspi_In)
	_RASPI4B_GPIOInit(raspi_DI3, raspi_In)
	_RASPI4B_GPIOInit(raspi_DI4, raspi_In)
}

func _RASPI4B_GPIOInit(Pin string, EnDir string) {
	//gpio export
	cmd := fmt.Sprintf("echo %s > /sys/class/gpio/export", Pin)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[RASPI4B_GPIOInit] error", err)
		return
	}
	// gpio set direction
	cmd = fmt.Sprintf("echo %s > /sys/class/gpio/gpio%s/direction", EnDir, Pin)
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[RASPI4B_GPIOInit] error", err)
	}
}

func RASPI4_GPIOSet(pin, value int) (bool, error) {
	cmd := fmt.Sprintf("echo %d > /sys/class/gpio/gpio%d/value", value, pin)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[RASPI4_GPIOSet] error", err)
		return false, err
	}
	v, e := RASPI4_GPIOGet(pin)
	if e != nil {
		return false, e
	}
	return v == value, nil
}

func RASPI4_GPIOGet(pin int) (int, error) {
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
