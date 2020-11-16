package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
)

func helloTask(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
	w.Write([]byte("hello"))
}
func hiTask(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hi\n")
	w.Write([]byte("hi"))
}

func main() {

	startHttpServer()
}

func startHttpServer() {
	//http := ModifierMiddleware
	serv := mux.NewRouter()
	http_port := "8090"

	//通过完整的path来匹配
	serv.HandleFunc("/api/login", helloTask)

	errser := http.ListenAndServe(":"+http_port, httpMiddleware(serv))
	if errser != nil {
		log.Println(errser)
	}

}

//获取当前运行目录
func GetCurrentPath() (dir string, err error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Println("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}
	dir = filepath.Dir(absPath)
	return dir, nil
}

//跨域处理的中间件 http服务所有的请求都会经过这里处理
func httpMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,Origin, X-Requested-With, Content-Type, Accept,common")

		h.ServeHTTP(w, r)

		if r.Method == "OPTIONS" {
			return
		}
	})
}
