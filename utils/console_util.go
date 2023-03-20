package utils

import (
	_ "embed"
	"fmt"
)

//go:embed banner.b
var banner string

//
// show banner
//
func ShowBanner() {

	fmt.Println("\n", banner)

}
