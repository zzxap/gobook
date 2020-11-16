//评测题目: 用程序实现一个函数，原型为：void printNumber(int N)
//内部创建三个线程A、B、C，三个线程交替输出1到N之间的数字，
//A线程输出1、4、7…，B线程2、5、8…，C线程输出3、6、9…，
//最终结果为：1、2、3、4、5...N

package main

import (
	"fmt"
	"sync/atomic"
)

var index int32

//开ABC三个channel 分别用于三个channel监听数据
var chA = make(chan int32)
var chB = make(chan int32)
var chC = make(chan int32)

//监听关闭
var chClose = make(chan int32)
var maxNum int32

type Values struct {
	m map[string]string
}

func (v Values) Get(key string) string {
	return v.m[key]
}
func main() {
	index = 0
	maxNum = 20
	printNumber(20)
}
func printNumber(N int) {
	//开三个协程做任务
	go printA(N)
	go printB(N)
	go printC(N)

	//给chA发数据通知它开始打印
	chA <- 1
	//等待
	<-chClose
	close(chA)
	close(chB)
	close(chC)
	close(chClose)
	fmt.Print(" finish")

}

func printA(N int) {
	//开始监听
	for {
		select {
		case v := <-chA:
			atomic.AddInt32(&index, 1)
			fmt.Print(v)
			fmt.Print("A,")
			chB <- v + 1 //通知B开始打印
			if v == maxNum {
				chClose <- 1
				return
			}
		case <-chClose:
			return
		}
	}

}
func printB(N int) {
	//B开始监听
	for {
		select {
		case v := <-chB:
			atomic.AddInt32(&index, 1)
			fmt.Print(v)
			fmt.Print("B,")
			chC <- v + 1 //输入chC 通知C开始打印

			if v == maxNum {
				chClose <- 1
				return
			}
		case <-chClose:
			return
		}
	}
}
func printC(N int) {
	//开始监听
	for {
		select {
		case v := <-chC:
			atomic.AddInt32(&index, 1)
			fmt.Print(v)
			fmt.Print("C,")
			chA <- v + 1 //通过channel通知A
			if v == maxNum {
				chClose <- 1
				return
			}
		case <-chClose:
			return
		}
	}
}
