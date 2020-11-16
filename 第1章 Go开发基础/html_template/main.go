package main

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Name   string
	Gender string
	Age    int
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	//加载模板 解析模板
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Println("Parse template failed, err%v", err)
		return
	}
	// 渲染字符串
	name := "mmm"
	//err = t.Execute(w, name)
	// 渲染结构体
	user := User{
		Name:   name,
		Gender: "男",
		Age:    23,
	}
	//err = t.Execute(w, user)
	// 渲染map
	m := map[string]interface{}{
		"name":   name,
		"gender": "男",
		"age":    24,
	}
	//err = t.Execute(w, m)
	carList := []string{
		"汽车",
		"火车",
		"货车",
	}
	//把对象传输到模板展示
	err = t.Execute(w, map[string]interface{}{
		"m":       m,
		"user":    user,
		"carList": carList,
	})
	if err != nil {
		log.Println("render template failed, err%v", err)
		return
	}
}
func main() {
	http.HandleFunc("/", sayHello)
	err := http.ListenAndServe(":8099", nil)
	if err != nil {
		log.Println("http server start failed,err:%v", err)
	}
}
