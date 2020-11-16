package main

import (
	"sync"
	"time"
)

var m *sync.RWMutex

func main() {
	m = new(sync.RWMutex)
	// 写的时候啥也不能干
	go write(1)
	go read(2)
	go write(3)

	time.Sleep(4 * time.Second)
}

func read(i int) {
	println(i, "读开始")
	m.RLock()
	println(i, "正在读...")
	time.Sleep(1 * time.Second)
	m.RUnlock()
	println(i, "读结束")
}

func write(i int) {
	println(i, "写开始")
	m.Lock()
	println(i, "正在写...")
	time.Sleep(1 * time.Second)
	m.Unlock()
	println(i, "写完毕")
}
