package main

import (
	"sync"
	"time"
)

var m *sync.RWMutex

func main() {
	m = new(sync.RWMutex)
	// 多个同时读
	go read(1)
	go read(2)
	time.Sleep(3 * time.Second)
}

func read(i int) {
	println(i, "开始读")
	m.RLock()
	println(i, "正在读...")
	time.Sleep(1 * time.Second)
	m.RUnlock()
	println(i, "读结束")
}
