//服务器端代码：
//这里暴露了一个RPC接口，一个HTTP接口

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
)

type Watcher string

func (w *Watcher) GetInfo(arg int, result *string) error {
	*result = "helloooooo"
	return nil
}

func main() {

	http.HandleFunc("/api", hello)

	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("监听失败，端口可能已经被占用")
	}
	fmt.Println("正在监听8000端口")
	http.Serve(l, nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><body>hello-123</body></html>")
}
