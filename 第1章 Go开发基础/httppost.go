package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	//发送请求
	//resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
	resp, err := http.PostForm("http://xx.com", url.Values{"q": {"github"}})
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

//demo 2
func post(url string, params map[string]string) {

	reqBody, err := json.Marshal(params)
	if err != nil {
		print(err)
	}
	resp, err := http.Post(url,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))

}
