/*

参考https://lanre.wtf/blog/2017/07/24/roundtripper-go/
1.缓存 responses，比如 app需要访问 Github api，获取 trending repos，这个数据变动不频繁，假设30分钟变动一次，你显然不希望每次都要点击api都要来请求Github api，解决这个问题的方法是实现这样的http.RoundTripper

有缓存时从缓存取出response数据
过期，数据通过重新请求api获取
2.根据需要设置http header, 一个容易想到的例子go-github一个Github的 api的go客户端。某些github api不需要认证，有些需要认证则需要提供自己的http client，比如 ghinstallation,下面是ghinstallation 的 RoundTrip 函数实现，设置 Authorization 头
image.png
3.限速(Rate limiting) 控制请求速率


*/

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {

	cachedTransport := newTransport()

	//Create a custom client so we can make use of our RoundTripper
	//If you make use of http.Get(), the default http client located at http.DefaultClient is used instead
	//Since we have special needs, we have to make use of our own http.RoundTripper implementation
	client := &http.Client{
		Transport: cachedTransport,
		Timeout:   time.Second * 5,
	}

	//Time to clear the cache store so we can make request to the original server
	cacheClearTicker := time.NewTicker(time.Second * 5)

	//Make a new request every second
	//This would help demonstrate if the response is actually coming from the real server or from the cache
	reqTicker := time.NewTicker(time.Second * 1)

	terminateChannel := make(chan os.Signal, 1)

	signal.Notify(terminateChannel, syscall.SIGTERM, syscall.SIGHUP)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000", strings.NewReader(""))

	if err != nil {
		log.Fatalf("An error occurred ... %v", err)
	}

	for {
		select {
		case <-cacheClearTicker.C:
			// Clear the cache so we can hit the original server
			cachedTransport.Clear()

		case <-terminateChannel:
			cacheClearTicker.Stop()
			reqTicker.Stop()
			return

		case <-reqTicker.C:

			resp, err := client.Do(req)

			if err != nil {
				log.Printf("An error occurred.... %v", err)
				continue
			}

			buf, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Printf("An error occurred.... %v", err)
				continue
			}

			fmt.Printf("The body of the response is \"%s\" \n\n", string(buf))
		}
	}
}

func cacheKey(r *http.Request) string {
	return r.URL.String()
}

type cacheTransport struct {
	data              map[string]string
	mu                sync.RWMutex
	originalTransport http.RoundTripper
}

func (c *cacheTransport) Set(r *http.Request, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("cache key:", cacheKey(r))
	fmt.Println("cache value:", value)
	c.data[cacheKey(r)] = value
}

func (c *cacheTransport) Get(r *http.Request) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.data[cacheKey(r)]; ok {
		return val, nil
	}

	return "", errors.New("key not found in cache")
}

// There be dragons!!!
func (c *cacheTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	// Check if we have the response cached..
	// If yes, we don't have to hit the server
	// We just return it as is from the cache store.
	if val, err := c.Get(r); err == nil {
		fmt.Println("Fetching the response from the cache")
		return cachedResponse([]byte(val), r)
	}

	// Ok, we don't have the response cached, the store was probably cleared.
	// Make the request to the server.
	resp, err := c.originalTransport.RoundTrip(r)

	if err != nil {
		return nil, err
	}

	// Get the body of the response so we can save it in the cache for the next request.
	buf, err := httputil.DumpResponse(resp, true)

	if err != nil {
		return nil, err
	}

	// Saving it to the cache store
	c.Set(r, string(buf))

	fmt.Println("Fetching the data from the real source")
	return resp, nil
}

func (c *cacheTransport) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]string)
	return nil
}

func cachedResponse(b []byte, r *http.Request) (*http.Response, error) {
	buf := bytes.NewBuffer(b)
	return http.ReadResponse(bufio.NewReader(buf), r)
}

func newTransport() *cacheTransport {
	return &cacheTransport{
		data:              make(map[string]string),
		originalTransport: http.DefaultTransport,
	}
}
