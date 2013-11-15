package webSocket

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var conn *websocket.Conn

func Start() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("../webSocket/index.html")
			t.Execute(w, nil)
		}
	})

	http.HandleFunc("/sendMsg", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if e2 := websocket.Message.Send(conn, strings.Join(r.Form["msg"], "")); e2 != nil {
			fmt.Println("send error")
			fmt.Println(e2)
		}
	})

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		for {
			var reply string
			if e1 := websocket.Message.Receive(ws, &reply); e1 != nil {
				fmt.Println("receive error")
				break
			}
			fmt.Println("receive:" + reply)
			msg := "receive msg :" + reply
			if e2 := websocket.Message.Send(ws, msg); e2 != nil {
				fmt.Println("send error")
				break
			}
			conn = ws
		}
	}))

	fmt.Println("listen...9999")
	//监听端口
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		log.Fatal(err)
	}
}
