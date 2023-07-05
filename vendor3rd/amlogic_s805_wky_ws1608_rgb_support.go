package vendor3rd

import archsupport "github.com/hootrhino/rulex/bspsupport"

// R string = "red"
// G string = "green"
// B string = "blue"
func AmlogicWKYS805_RGBSet(pin string, value int) (bool, error) {
	return archsupport.AmlogicWKYS805_RGBSet(pin, value)
}

func AmlogicWKYS805_RGBGet(pin string) (int, error) {
	return archsupport.AmlogicWKYS805_RGBGet(pin)
}
