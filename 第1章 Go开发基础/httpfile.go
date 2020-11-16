package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

	//静态文件路由
	//把xx.png文件放在file目录下 就可以通过 http://127.0.0.1:8090/file/xx.png 下载了
	curdir, _ := GetCurrentPath()
	PthSep := string(os.PathSeparator)
	filePath := curdir + PthSep + "file"
	router.PathPrefix("/file/").Handler(http.StripPrefix("/file/", http.FileServer(http.Dir(filePath))))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

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
