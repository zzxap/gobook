package main

import (
	"bufio"
	"fmt"
	"net"

	//"os"
	"strconv"
	"strings"
)

var count = 0

func handleConnection(c net.Conn) {
	fmt.Print(".")
	for {
		//读取客户端数据
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		fmt.Println(temp)
		counter := strconv.Itoa(count) + "\n"
		//给客户端发送数据
		c.Write([]byte(string(counter)))
	}
	c.Close()
}

func main() {

	PORT := ":8082"
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		//每多一个连接count就加一
		go handleConnection(c)
		count++
	}
}
