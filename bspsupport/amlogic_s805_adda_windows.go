package archsupport

/*
*
* Windows 不支持 特定实现
*
 */

func init() {
	//fmt.Println("WKYS805 RGB GPIO not support Windows")
}

/*
*
* Windows下我们不实现 但是留下接口防止未来扩展的时候改代码
*
 */

func AmlogicWKYS805_RGBSet(pin string, value int) (bool, error) {
	return false, errArchNotSupport
}

func AmlogicWKYS805_RGBGet(pin string) (int, error) {

	return -1, errArchNotSupport
}
