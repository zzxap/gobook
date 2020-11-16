package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8000"
	CONN_TYPE = "tcp"
)

func main() {
	// 监听传入的连接。
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// 当应用程序关闭时，关闭监听器。
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// 监听传入的连接。
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// 在新的goroutine中处理连接。
		go handleRequest(conn)
	}
}

//处理传入的请求。
func handleRequest(conn net.Conn) {
	// 创建一个缓冲区以保存传入的数据。
	buf := make([]byte, 1024)
	//将传入的连接读入缓冲区。
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("receive: ", string(buf))

	// 将回复发送给与我们联系的人。
	conn.Write([]byte("Server :Message received."))
	// 完成连接后，关闭连接。
	conn.Close()
}
