/*
包gorilla/mux实现了一个请求路由器和调度程序，用于将传入的请求与它们各自的处理程序进行匹配。

名称mux代表“ HTTP请求多路复用器”。与standard一样http.ServeMux，mux.Router
将传入的请求与已注册的路由列表进行匹配，
并为与URL或其他条件匹配的路由调用处理程序。主要特点是：

它实现了http.Handler接口，因此与标准兼容http.ServeMux。
可以基于URL主机，路径，路径前缀，方案，标头和查询值，HTTP方法或使用自定义匹配器来匹配请求。
URL主机，路径和查询值可以具有带可选正则表达式的变量。
可以构建或“反转”已注册的URL，这有助于维护对资源的引用。
路由可用作子路由：仅在父路由匹配时才测试嵌套路由。这对于定义具有共同条件（例如主机，路径前缀或其他重复属性）
的路由组很有用。另外，这可以优化请求匹配。
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func helloTask(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
	w.Write([]byte("hello"))
}

func main() {
	startHttpServer()
}

func startHttpServer() {
	//http := ModifierMiddleware
	router := mux.NewRouter()

	//通过完整的path来匹配
	router.HandleFunc("/api/hello", helloTask)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8090",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
