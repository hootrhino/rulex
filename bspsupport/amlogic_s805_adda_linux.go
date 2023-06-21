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

/*
玩客云WS1608有一个RGB LED，其系统内部已经直接支持进设备树：

	R -> /sys/class/leds/onecloud\:red\:alive/brightness
	G -> /sys/class/leds/onecloud\:green\:alive/brightness
	B -> /sys/class/leds/onecloud\:blue\:alive/brightness
*/
const (
	AmlogicWKYS805_R string = "red"
	AmlogicWKYS805_G string = "green"
	AmlogicWKYS805_B string = "blue"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "WKYS805" {
		_AmlogicWKYS805_RGBAllInit()
	}
}

/*
explain:init all gpio
*/
func _AmlogicWKYS805_RGBAllInit() int {
	_AmlogicWKYS805_RGBInit(AmlogicWKYS805_R, 0)
	_AmlogicWKYS805_RGBInit(AmlogicWKYS805_G, 0)
	_AmlogicWKYS805_RGBInit(AmlogicWKYS805_B, 0)
	// 返回值无用
	return 1
}

func _AmlogicWKYS805_RGBInit(pin string, value int) {
	AmlogicWKYS805_RGBSet(pin, value)
}

func AmlogicWKYS805_RGBSet(pin string, value int) (bool, error) {
	cmd := fmt.Sprintf("echo %d > /sys/class/leds/onecloud\\:%s\\:alive/brightness", value, pin)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("[AmlogicWKYS805_RGBSet] error", err)
		return false, err
	}
	v, e := AmlogicWKYS805_RGBGet(pin)
	if e != nil {
		return false, e
	}
	return v == value, nil
}

func AmlogicWKYS805_RGBGet(pin string) (int, error) {
	cmd := fmt.Sprintf("cat /sys/class/leds/onecloud\\:%s\\:alive/brightness", pin)
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
