package subpackage

import (
	"fmt"
)

func init() {
	fmt.Println("this is in subpk")
}

//同一个包下不能有相同方法
// func M2() {
// 	fmt.Println("this is in method M2 in subpk")
// }
