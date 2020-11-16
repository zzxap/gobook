package main

import (
	"fmt"
	"time"
)

// goroutine 1
func Bname() {
	arr1 := [4]string{"aa", "bb", "cc", "dd"}
	for t1 := 0; t1 <= 3; t1++ {

		time.Sleep(150 * time.Millisecond)
		fmt.Printf("%s\n", arr1[t1])
	}
}

// goroutine 2
func Bid() {
	arr2 := [4]int{11, 22, 33, 44}
	for t2 := 0; t2 <= 3; t2++ {
		time.Sleep(150 * time.Millisecond)
		fmt.Printf("%d\n", arr2[t2])
	}
}

// Main function
func main() {
	fmt.Println("!...主协程开始...!")
	// 创建运行 Goroutine 1
	go Bname()
	// 创建运行 Goroutine 2
	go Bid()
	time.Sleep(3500 * time.Millisecond)
	fmt.Println("\n!...主协程结束...!")
}

/*
运行结果
!...主协程开始...!
11
aa
22
bb
33
cc
44
dd

!...主协程结束...!

*/
