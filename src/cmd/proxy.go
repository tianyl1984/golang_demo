package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"net/http"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	err := http.ListenAndServe(":8888", proxy)
	if err != nil {
		fmt.Println(err)
	}
}
