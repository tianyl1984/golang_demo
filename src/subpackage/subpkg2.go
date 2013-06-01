package subpackage

import (
	"fmt"
)

func init() {
	fmt.Println("this is in subpkg2")
}

func M2() {
	fmt.Println("this is in method M2")
}
