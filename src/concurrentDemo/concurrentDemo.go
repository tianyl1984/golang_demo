package concurrentDemo

import (
	"fmt"
	"runtime"
)

func StartDemo() {
	m1()
}

func m1() {
	go func() {
		for i := 0; i < 5; i++ {
			runtime.Gosched() //cpu把时间片让给其他goroutines
			fmt.Println("Hello")
		}
	}()
	fmt.Println("World") //当前goroutines执行
}
