package osDemo

import (
	"fmt"
	"os"
	"os/exec"
)

func EnvDemo() {
	os.Setenv("appId", "hello")        //设置环境变量
	fmt.Println(os.Getenv("appId"))    //读取环境变量
	for _, env := range os.Environ() { //遍历环境变量
		fmt.Println(env)
	}
}

func ExeDemo() {
	cmd := exec.Command("ipconfig")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Command Error!", err.Error())
		return
	}
	fmt.Println(string(out))
}
