package archsupport

/*
*
* Windows下我们不实现 但是留下接口防止未来扩展的时候改代码
*
 */

func RASPI4_GPIOSet(pin, value int) (bool, error) {
	return false, nil
}
func RASPI4_GPIOGet(pin int) (int, error) {
	return 0, nil
}
