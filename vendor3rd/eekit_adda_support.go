package vendor3rd

import (
	"github.com/hootrhino/rulex/archsupport"
)

/*
*
* 跨平台支持
*
 */
func EEKIT_GPIOSet(pin, value int) (bool, error) {
	return archsupport.EEKIT_GPIOSet(pin, value)
}
func EEKIT_GPIOGet(pin int) (int, error) {
	return archsupport.EEKIT_GPIOGet(pin)
}
