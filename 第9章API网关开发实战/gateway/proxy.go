/*
https://blog.charmes.net/post/reverse-proxy-go/

正向代理隐藏了真实客户端向服务器发送请求，反向代理隐藏了真实服务器向客户端提供服务。
这段代码比较直观，只包含了最核心的代码逻辑，完全按照最上面的代理图例进行组织。一共分成几个步骤：
代理接收到客户端的请求，复制了原来的请求对象，并根据数据配置新请求的各种参数（添加上 X-Forward-For 头部等）
把新请求发送到服务器端，并接收到服务器端返回的响应
代理服务器对响应做一些处理，然后返回给客户端

*/

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
	transport := http.DefaultTransport
	// step 1
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	// step 2
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	// step 3
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("0.0.0.0:8080", nil)
}
