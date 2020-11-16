package main

import (
	"fmt"
	//"strconv"
	//"sync"
	"time"
)

func main() {
	//测试
	testselect()
	//testchanel()
}

func testChannel22() {
	c := make(chan int) // Allocate a channel.

	// Start the sort in a goroutine; when it completes, signal on the channel.
	go func() {
		//list.Sort()
		c <- 1 // Send a signal; value does not matter.
	}()

	//doSomethingForAWhile()
	<-c
}

func testchanel() {
	//channel 的读写操作
	//channel 一定要初始化后才能进行读写操作，否则会永久阻塞。
	ch := make(chan int, 10)
	// x写入channel
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	// 读 channel

	for x := range ch {
		fmt.Println(x)
	}

	x, ok := <-ch

	fmt.Println(x, ok)

	//关闭 channel
	close(ch)

	/*
		channel 的关闭，你需要注意以下事项:

		关闭一个未初始化(nil) 的 channel 会产生 panic
		重复关闭同一个 channel 会产生 panic
		向一个已关闭的 channel 中发送消息会产生 panic
		从已关闭的 channel 读取消息不会产生 panic，
		且能读出 channel 中还未被读取的消息，若消息均已读出，则会读到类型的零值。
		从一个已关闭的 channel 中读取消息永远不会阻塞，并且会返回一个为 false 的 ok-idiom，
		可以用它来判断 channel 是否关闭
		关闭 channel 会产生一个广播机制，所有向 channel 读取消息的 goroutine 都会收到消息
	*/
}
func fprint(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}
func testGoroutine() {

	fprint("direct")

	go fprint("goroutine")

	go func(msg string) {
		fmt.Println(msg)
	}("going")

	time.Sleep(time.Second)
	fmt.Println("done")

}

//https://gobyexample.com/
func testselect() {

	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}

}
