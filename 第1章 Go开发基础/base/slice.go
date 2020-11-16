package main

import (
	"fmt"
)

func main() {
	//定义一个无初始长度的切片
	a := []string{}
	for i := 0; i < 6; i++ {
		a = append(a, "11111")
		//Slice的容量
		fmt.Println(cap(a))
		//Slice的长度
		fmt.Println(len(a))
		fmt.Println("-------")
	}
}
