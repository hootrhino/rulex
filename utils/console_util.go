package utils

import (
	_ "embed"
	"fmt"
)

//go:embed banner.txt
var banner string

//
// show banner
//
func ShowBanner() {

	fmt.Println("\n", banner)

}
