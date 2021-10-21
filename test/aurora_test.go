package test

import (
	"fmt"
	"testing"

	. "github.com/logrusorgru/aurora"
)

func Test_aurora(t *testing.T) {
	fmt.Println("Hello,", Magenta("Aurora"))
	fmt.Println(Bold(Cyan("Cya!")))
}
