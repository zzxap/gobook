package main

import (
	"fmt"
	"time"
)

func selectExample() {
	//For our example we’ll select across two channels.

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

/*
received one
received two
*/

func sele() {
	/*
		从官方文档看一下有关select的描述：

		A "select" statement chooses which of a set of possible send or receive operations will proceed. It looks similar to a "switch" statement but with the cases all referring to communication operations.
		一个select语句用来选择哪个case中的发送或接收操作可以被立即执行。它类似于switch语句，但是它的case涉及到channel有关的I/O操作。

		或者换一种说法，select就是用来监听和channel有关的IO操作，当 IO 操作发生时，触发相应的动作。
		如果有一个或多个IO操作可以完成，则Go运行时系统会随机的选择一个执行，否则的话，如果有default分支，则执行default分支语句，如果连default都没有，则select语句会一直阻塞，直到至少有一个IO操作可以进行

		select 可以同时监听多个 channel 的写入或读取
		执行 select 时，若只有一个 case 通过(不阻塞)，则执行这个 case 块
		若有多个 case 通过，则随机挑选一个 case 执行
		若所有 case 均阻塞，且定义了 default 模块，则执行 default 模块。
		若未定义 default 模块，则 select 语句阻塞，直到有 case 被唤醒。
		使用 break 会跳出 select 块。


	*/
	//select基本用法

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}

}
