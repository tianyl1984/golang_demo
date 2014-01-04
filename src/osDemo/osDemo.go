package osDemo

import (
	"fmt"
	"os"
)

func EnvDemo() {
	os.Setenv("appId", "hello")        //设置环境变量
	fmt.Println(os.Getenv("appId"))    //读取环境变量
	for _, env := range os.Environ() { //遍历环境变量
		fmt.Println(env)
	}
}
