package osDemo

import (
	"flag"
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

func ExeDemo() { //执行外部命令
	cmd := exec.Command("ipconfig")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Command Error!", err.Error())
		return
	}
	fmt.Println(string(out))
}

func ArgsDemo() { //读取参数
	fmt.Println(os.Args)

	//-url=http://www.xxx.com
	url := flag.String("url", "http://www.baidu.com", "网址")
	flag.Parse()
	fmt.Println(*url)
}
