package fileDemo

import (
	"fmt"
	"os"
)

func StartDemo() {
	m1()
}

func m1() {
	//文件夹操作
	path := "aaa/bbb"
	fi, err := os.Stat("aaa/bbb")
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("文件夹存在")
		} else {
			fmt.Println("文件夹不存在")
			os.MkdirAll(path, 0777)
			fmt.Println("已创建文件夹")
		}
	} else {
		if fi.IsDir() {
			fmt.Println("文件夹存在！！")
		} else {
			fmt.Println(path, "不是目录！！")
		}
		os.RemoveAll(path)
		fmt.Println("已删除", path, "！！")
	}
	f, err := os.Create("aaa.txt") //已存在会覆盖
	checkError(err)
	fmt.Println("已创建文件：", f.Name())
	f.WriteString("测试文件！！！！！！")
	defer f.Close()
	f2, err := os.Open("login.gtpl")
	defer f2.Close()
	checkError(err)
	buf := make([]byte, 1024)
	for {
		n, _ := f2.Read(buf)
		if 0 == n {
			break
		}
		fmt.Println(string(buf[:n]))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
