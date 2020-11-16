package main

import (
	"fmt"
	//"strconv"
	"sync"
	"time"
)

func testAsync() {
	fmt.Println("aaa \n")
	//开启一个协程去执行任务
	go asyncFunc("bbb1 \n")
	fmt.Println("ccc \n")

}
func testAsyncSleep() {
	fmt.Println("aaa \n")
	//开启一个协程去执行任务
	go asyncFunc("bbb2 \n")
	fmt.Println("ccc \n")
	//在结束前等待一下
	time.Sleep(1 * time.Second)
}

func asyncFunc(str string) {
	fmt.Println(str)
}

//使用WaitGroups控制并发，等待返回
func testAsyncWait() {
	fmt.Println("aaa")

	var waitgroup sync.WaitGroup
	//开启协程执行一个任务 就add(1) 执行完就Add(-1)或done
	waitgroup.Add(10)
	//异步执行10个任务
	for i := 0; i < 10; i++ {

		go func(index int) {
			fmt.Println(index)
			//任务执行完毕 协程个数减1
			waitgroup.Done()
			//waitgroup.Add(-1)
		}(i)
	}
	//这里一直在等待，waitgroup里的任务数量清零
	waitgroup.Wait()

	fmt.Println("ccc")
	/*
		输出 其中1-10 是随机的乱序的，每次执行的次序都不一样
		aaa
		0
		1
		9
		6
		7
		5
		8
		2
		4
		3
		ccc
	*/

}

//如果WaitGroup 执行超时 waitgroup.Wait() 一直在等待怎么办呢？可以加一个超时控制
/*
首先，我们来从官方文档看一下有关select的描述：

A "select" statement chooses which of a set of possible send or receive
 operations will proceed. It looks similar to a "switch" statement but with
the cases all referring to communication operations.
一个select语句用来选择哪个case中的发送或接收操作可以被立即执行。它类似于switch语句，
但是它的case涉及到channel有关的I/O操作。

*/
func testWaitGroupTimeOut() {
	var w = sync.WaitGroup{}
	var ch = make(chan bool)
	w.Add(2)
	//执行任务1
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("等2秒")
		w.Done()
	}()
	//执行任务2
	go func() {
		time.Sleep(time.Second * 6)
		fmt.Println("等6秒")
		w.Done()
	}()
	go func() {
		w.Wait()
		//执行完毕，向ch写入数据
		ch <- false
	}()
	//select就是用来监听和channel有关的IO操作，当 IO 操作发生时，触发相应的动作。
	select {
	case <-time.After(time.Second * 5):
		fmt.Println("超时了")
	case <-ch:
		// 如果成功向ch写入数据，则进行该case处理语句
		fmt.Println("结束了")
	}
}

func main() {
	//测试
	testWaitGroupTimeOut()

	/*
			   输出
			    aaa  ccc
		        并没有输出bbb，原因是主程序在协程执行之前就已经退出了。
				如果要等待bbb输出，必须等待足够的事件等待asyncFunc执行完毕
				执行testAsyncWait 将会看到输出
				 aaa  ccc  bbb1
	*/

}
