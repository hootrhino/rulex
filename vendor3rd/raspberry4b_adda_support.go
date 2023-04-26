package vendor3rd

import (
	"github.com/hootrhino/rulex/archsupport"
)

/*
*
* 跨平台支持
*
 */

func RASPI4_GPIOSet(pin, value int) (bool, error) {
	return archsupport.RASPI4_GPIOSet(pin, value)
}
func RASPI4_GPIOGet(pin int) (int, error) {
	return archsupport.RASPI4_GPIOGet(pin)
}
