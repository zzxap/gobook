//http client 提供更复杂强大的自定义请求

package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	data := url.Values{}

	data.Add("username", "aaa")
	data.Add("password", "111")
	urls := "http://127.0.0.1:8094/login"
	body, err := myhttpRequest(urls, data)
	if err != nil {

		log.Println(err.Error())

	}
	log.Println(string(body))
}

var transport *http.Transport
var client *http.Client

func myhttpRequest(url string, params url.Values) (body []byte, err error) {

	if transport == nil {
		if strings.Contains(url, "https") {

			transport = &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			}
		} else {
			transport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			}
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}

	}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Println("Error Occured. %+v", err)
		return nil, err
	}
	//("Authorization", " Bearer " + authorization);
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, errr := client.Do(req)
	//res, errr := client.Post(url, "application/x-www-form-urlencoded", nil) //strings.NewReader("name=cjb")
	if errr != nil {
		log.Println("client.Post error")
		log.Println(errr)
		log.Println(url)
		return nil, errr
	}

	bodyy, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("client.Post read error")
		log.Println(err)
		return nil, err
	}
	res.Body.Close()
	return bodyy, err

}
