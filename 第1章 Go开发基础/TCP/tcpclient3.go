package main

import (
	"fmt"
	"net"
	"time"
)

var conn net.Conn
var err error

func main() {

	CONNECT := "127.0.0.1:8082"
	conn, err = net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 20; i++ {
		fmt.Println("write")
		time.Sleep(time.Duration(2) * time.Second)
		//conn.Write([]byte("sdfsdfsdf"))
		fmt.Fprintf(conn, "32233232\n")
	}

}
