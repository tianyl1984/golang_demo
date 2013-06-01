package subpackage

import (
	"fmt"
)

const a string = "this is in subpkg"

func init() {
	fmt.Println(a)
}

func M() {
	fmt.Println("this is in method M")
}
