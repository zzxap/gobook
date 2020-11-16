/*

Go 1.3 的sync包中加入一个新特性：Pool。
这个类设计的目的是用来保存和复用临时对象，以减少内存分配，降低CG压力。

sync.Pool是一个可以存或取的临时对象集合
sync.Pool可以安全被多个线程同时使用，保证线程安全
注意、注意、注意，sync.Pool中保存的任何项都可能随时不做通知的释放掉，
所以不适合用于像socket长连接或数据库连接池。
sync.Pool主要用途是增加临时对象的重用率，减少GC负担。

*/
package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Pool for our struct A
var pool *sync.Pool

// A dummy struct with a member
type Person struct {
	Name string
}

// Func to init pool
func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			//fmt.Println("返回A")
			return new(Person)
		},
	}
}

// Main func
func main() {
	// Initializing pool
	initPool()
	//sync.Pool中保存的任何项都可能随时不做通知的释放掉，所以不适合用于像socket长连接或数据库连接池
	// Get hold of instance one
	//需要新的结构体的时候，尝试去pool中取，而不是重新生成，重复10000次仅需要节省大量的时间。
	//这样简单的操作，却节约了99.65%的时间，也节约了各方面的资源。
	//最重要的是它可以有效减少GC CPU和GC Pause。
	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		one := pool.Get().(*Person)
		one.Name = "girl" + strconv.Itoa(i)
		//fmt.Printf("one.Name = %s\n", one.Name)
		//使用后提交回实例
		pool.Put(one)

	}
	// 现在，同一实例可被另一个例程使用，而无需再次分配它
	fmt.Println("花费时间1:", time.Since(startTime))

	startTime = time.Now()
	for i := 0; i < 10000; i++ {

		p := Person{Name: "girl" + strconv.Itoa(i)}
		pool.Put(p)
	}

	fmt.Println("花费时间2:", time.Since(startTime))

}

/*
使用pool会快一点

花费时间1: 1.0079ms
花费时间2: 1.9395ms
*/
