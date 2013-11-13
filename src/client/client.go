package client

import (
	"fmt"
	// "io/ioutil"
	"net"
	// "os"
)

func init() {
	fmt.Println("client start")
}

func m1() {
	result := make([]byte, 128)
	for {
		fmt.Println("input cmd:")
		var str string
		_, err := fmt.Scanf("%s", &str)
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		// fmt.Println("you input:", str)
		if str == "exit" {
			break
		}
		tcpAdd, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
		checkError(err)
		conn, err := net.DialTCP("tcp", nil, tcpAdd) //设置local address为nil
		checkError(err)
		_, err = conn.Write([]byte(str))
		checkError(err)
		//result, err := ioutil.ReadAll(conn)
		read_len, err := conn.Read(result)
		if err != nil {
			fmt.Println(err)
			break
		}
		checkError(err)
		//fmt.Println("response:", string(result))
		fmt.Println(string(result[:read_len]))
		conn.Close()
		result = make([]byte, 128)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func m2() {
	str := "123"
	result := make([]byte, 128)
	tcpAdd, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6000")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAdd) //设置local address为nil
	checkError(err)
	_, err = conn.Write([]byte(str))
	checkError(err)
	//result, err := ioutil.ReadAll(conn)
	read_len, err := conn.Read(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	checkError(err)
	//fmt.Println("response:", string(result))
	fmt.Println(string(result[:read_len]))
	conn.Close()
}

func Start() {
	//m1()
	m2()
}
