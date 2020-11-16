package main

import (
	"fmt"
)

func main() {
	// map
	m := make(map[int]string)
	m[0] = "a"
	m[1] = "b"
	changeMap(m)
	fmt.Printf("map:%+v", m) //输出 map:map[0:aaa 1:b]
	fmt.Println()

	//array
	var a = [2]string{"a", "b"}
	changeArray(a)
	fmt.Printf("array:%+v", a) //输出array:[a b]
	fmt.Println()

	//slice
	var s = []string{"a", "b"}
	changeSlice(s)
	fmt.Printf("slice:%+v", s) //输出slice:[aaa b]
}

func changeMap(m map[int]string) {
	m[0] = "aaa"
}

func changeArray(a [2]string) {
	a[0] = "aaa"
}

func changeSlice(s []string) {
	s[0] = "aaa"
}
