package fakeServer

import (
	"fmt"
	"net/http"
)

func Start() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><head><meta http-equiv='Content-Type' content='text/html; charset=gb2312' /><title></title></head><body><input type='button' value=启动 onclick=\"var a = window.open('http://www.baidu.com');a.close();\"></body></html>")
	})
	//<iframe height='125' width='835' frameborder='no' scrolling='no' src= 'javascript:void(0)'></iframe>
	//http://www.ueads.net/code/adview_pic6.php?r=1&c=7&w=835&h=125&b=0080ff&s=004080&bg=FFFFFF&p=808080&u=1925&at=p0&tt=t1
	fmt.Println("start...80")
	//监听端口
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
}
