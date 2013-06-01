package main

import (
	"fmt"
	"unicode/utf8"
)

func Q2() {
	fmt.Println("---------1----------")
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	fmt.Println("---------2----------")
	j := -1
lable:
	j++
	if j < 10 {
		fmt.Println(j)
		goto lable
	}
	fmt.Println("---------3----------")
	arr := [...]int{1, 2, 3, 4, 5, 6}
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}

func Q3() {
	for i := 1; i <= 100; i++ {
		switch {
		case i%15 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}

func Q4() {
	fmt.Println("---------1----------")
	for i := 1; i <= 10; i++ {
		var s string = ""
		for j := 0; j < i; j++ {
			s += "A"
		}
		fmt.Println(s)
	}
	fmt.Println("---------2----------")
	s := "asdfj找那个sd"
	fmt.Println("字节数：", len(s))
	fmt.Println("字符数：", utf8.RuneCountInString(s))
	fmt.Println("---------3----------")
	fmt.Println(s)
	//将字符串转为rune（int32）的silce类型，每个rune保存unicode编码的指针
	runeslice := []rune(s)
	runeslice[3] = 65
	runeslice[4] = 66
	runeslice[5] = 67
	s = string(runeslice)
	fmt.Println(s)
	s = "abcd中国"
	runeslice = []rune(s)
	for i := 0; i < len(runeslice)/2; i++ {
		runeslice[i], runeslice[len(runeslice)-1-i] = runeslice[len(runeslice)-1-i], runeslice[i]
	}
	s = string(runeslice)
	fmt.Println(s)
}

func Q5() {
	slice64 := []float64{1.2, 2.3, 4.5, 6.7}
	sum := 0.0
	for i := range slice64 {
		sum += slice64[i]
	}
	result := sum / float64(len(slice64))
	fmt.Println(result)
}

func Q7(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func Q9() {
	fmt.Println("---------1----------")

	fmt.Println("---------2----------")

}

func main() {
	// Q2()
	// Q3()
	// Q4()
	// Q5()
	fmt.Println(Q7(4, 2))
}
