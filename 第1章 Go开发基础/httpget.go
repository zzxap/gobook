package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	//发送请求

	resp, err := http.Get("https://www.xxxxx.com")

	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	//读取回复内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Print(string(body))
}
