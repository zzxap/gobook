package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	//以下变量定义只是方便演示,如果有自己的命名规则可以按照自己的来，不必纠结
	//整形  统一加冒号 这样a被自动识别为int类型
	a := 10
	//也可以这样定义 var a int=10 不过这样不够简洁，你觉得呢？

	//浮点型 这样b被自动识别为float类型 默认为 float64
	b := 3.14

	//同时定义多个变量 c是整形 d是字符串 e是float
	c, d, e := 1, "hello", 3.14

	//字符
	//只能被单引号包裹，不能用双引号，双引号包裹的是字符串
	f := 'a'                    //或 var f byte = 'a'
	fmt.Printf("%d %T\n", f, f) // 输出 97 uint8
	//字符串
	g := "xyz"

	//布尔型
	h := false
	fmt.Printf("a = %v, b = %v, c= %v,d = %v,e = %v,f = %v,g = %v,h = %v ", a, b, c, d, e, f, g, h)

	//基本数据类型默认值
	//在Golang中，数据类型都有一个默认值，当程序员没有赋值时，就会保留默认值，
	//在Golang中，默认值也叫做零值。
	var as int
	var bs float32
	var isTrue bool
	var str string

	//这里的%v,表示按照变量的值输出
	fmt.Printf("as = %v, bs = %v, isTrue = %v, str = %v", as, bs, isTrue, str)
	//打印结果as = 0, bs = 0, isTrue = false, str =
	array()
	multiArray()
	slice()
	maptest()
	syncmaptest()

}

func array() {
	//int数组
	var a [3]int             // 定义三个整数的数组
	fmt.Println(a[0])        // 打印第一个元素
	fmt.Println(a[len(a)-1]) // 打印最后一个元素
	// 打印索引和元素
	for i, v := range a {
		fmt.Printf("%d %d\n", i, v)
	}
	// 仅打印元素
	for _, v := range a {
		fmt.Printf("%d\n", v)

	}
	//字符串数组
	var team [3]string
	team[0] = "aa"
	team[1] = "bb"
	team[2] = "cc"
	//或者 team := [...]string{"aa", "bb", "cc"}

	for k, v := range team {
		fmt.Println(k, v)
	}
	//定义一个有三个元素的整型数组
	g := [3]int{12, 78, 50}
	fmt.Println(g)

	//声明了一个长度为 3 的数组，但是只提供了一个初值 12。剩下的两个元素被自动赋值为 0
	b := [3]int{12}
	fmt.Println(b)
	//输出为：[12 0 0]。

	//在声明数组时你可以忽略数组的长度并用 ... 代替，让编译器自动推导数组的长度
	c := [...]int{12, 78, 50}
	fmt.Println(c)

	//数组的遍历
	d := [...]float64{22.7, 23.8, 56, 68, 78}
	for i := 0; i < len(d); i++ {
		fmt.Printf("%d th element of a is %.2f\n", i, d[i])
	}
	//Go 提供了一个更简单，更简洁的遍历数组的方法：
	//使用 range for。range 返回数组的索引和索引对应的值。
	for i, v := range d { //range returns both the index and value
		fmt.Printf("%d the element of a is %.2f\n", i, v)
	}
	//输出如下
	/*
		0 the element of a is 22.70
		1 the element of a is 23.80
		2 the element of a is 56.00
		3 the element of a is 68.00
		4 the element of a is 78.00
	*/

	//如果你只想访问数组元素而不需要访问数组索引，则可以通过空标识符来代替索引变量
	for _, v := range d {
		fmt.Printf(" %.2f\n", v)
	}
}
func getSomethine() (string, int, float) {
	return "aa", 1, 1.22
}

//多维数组
func multiArray() {

	a := [3][2]string{
		{"a1", "a2"},
		{"b1", "b2"},
		{"c1", "c2"},
	}
	printarray(a)

	//初始化二维数组的另一种方式
	var b [3][2]string
	b[0][0] = "a1"
	b[0][1] = "a2"
	b[1][0] = "b1"
	b[1][1] = "b2"
	b[2][0] = "c1"
	b[2][1] = "c2"

	printarray(b)
}

//遍历多维数组
func printarray(a [3][2]string) {
	for _, v1 := range a {
		for _, v2 := range v1 {
			fmt.Printf("%s ", v2)
		}
		fmt.Printf("\n")
	}
}

/*
尽管数组看起来足够灵活，但是数组的长度是固定的，没办法动态增加数组的长度。
而切片却没有这个限制，实际上在 Go 中，切片比数组更为常见
一个数组不能动态改变长度。不要担心这个限制，因为切片（slices）可以弥补这个不足。
*/

func slice() {
	//创建了一个长度为 3 的 int 数组，并返回一个切片给 c
	c := []int{6, 7, 8}
	//修改切片
	c[1] = 70
	//c变成 [6 70 8]

	numa := [3]int{78, 79, 80}
	nums1 := numa[:]
	for _, v1 := range nums1 {
		fmt.Printf("%d ", v1)
	}

	//numa[:] 中缺少了开始和结束的索引值，
	//这种情况下开始和结束的索引值默认为 0 和len(numa) 表示整个数组

	//用 make 创建切片
	//内置函数 func make([]T, len, cap) []T 可以用来创建切片，该函数接受长度和容量作为参数
	si := make([]int, 5, 5)
	//用 make 创建的切片的元素值默认为 0 值。上面的程序输出为：[0 0 0 0 0]。
	for _, v1 := range si {
		fmt.Printf("%d ", v1)
	}
	//追加元素到切片
	cars := []string{"a", "b", "c"}
	cars = append(cars, "d")
	//cars变成[a b c d]
	cars = append(cars, "e", "f", "g", "h")
	//cars变成[a b c d e f g h]

	//合并切片 可以使用 ... 操作符将一个切片追加到另一个切片末尾
	veggies := []string{"potatoes", "tomatoes", "brinjal"}
	fruits := []string{"oranges", "apples"}
	food := append(veggies, fruits...)
	//变成 [potatoes tomatoes brinjal oranges apples]

	for _, v1 := range food {
		fmt.Printf("%s ", v1)
	}

	//删除切片第1个元素
	index := 1
	food = append(food[:index], food[index+1:]...)

	for _, v1 := range food {
		fmt.Printf("%s ", v1)
	}
	//清空
	food = append([]string{})

	//多维切片
	pls := [][]string{
		{"potatoes", "tomatoes"},
		{"brinjal"},
		{"oranges", "apples"},
	}
	for _, v1 := range pls {
		for _, v2 := range v1 {
			fmt.Printf("%s ", v2)
		}
		fmt.Printf("\n")
	}
}

//关于切片的内存回收
func sliceRecycle() {
	cars := []string{"ford", "toyota", "ds", "honda", "suzuki"}
	neededcars := cars[:len(cars)-2]
	//切片保留对底层数组的引用。只要切片存在于内存中，此时数组就不能被垃圾回收。

	carsCpy := make([]string, len(neededcars))
	copy(carsCpy, neededcars)
	//使用 copy 函数 func copy(dst, src []T) int 来创建该切片的一个拷贝。
	//这样我们就可以使用这个新的切片，原来的数组可以被垃圾回收。
	//现在数组cars 可以被垃圾回收，因为 neededcars 不再被引用。
	//return carsCpy
}
func maptest() {
	// 先声明map
	var m1 map[string]string
	//声明后指向的是nil，所以千万别直接声明后就使用
	//Golang中，map是引用类型，如切片一样，
	//再使用make函数创建一个非nil的map，nil map不能赋值
	m1 = make(map[string]string)
	// 最后给已声明的map赋值
	m1["a"] = "a"
	m1["b"] = "b"

	// 直接创建
	m2 := make(map[string]string)
	// 然后赋值
	m2["a"] = "a"
	m2["b"] = "b"

	// 初始化 + 赋值一体化
	m3 := map[string]string{
		"a": "a1",
		"b": "b1",
	}
	fmt.Println(m3)
	// 查找键值是否存在
	if v, ok := m1["a"]; ok {
		fmt.Println(v)
	} else {
		fmt.Println("键值不存在")
	}

	// 遍历map
	for k, v := range m1 {
		fmt.Println(k, v)
	}
}

//在并发中使用map
func syncmaptest() {
	c := make(map[string]string)
	//通过WaitGroup控制并发阻塞线程
	wg := sync.WaitGroup{}
	//同步互斥锁
	var lock sync.Mutex
	for i := 0; i < 10; i++ {
		//wg.Add(1) 和wg.Done() 要对应否则wg.Wait()会一直等待
		wg.Add(1)
		go func(n int) {
			k, v := strconv.Itoa(n), strconv.Itoa(n)
			//锁住开始修改
			lock.Lock()
			c[k] = v
			//修改完毕，解锁
			lock.Unlock()
			//wg.Done() 告诉WaitGroup，我完成一个任务了
			wg.Done()
		}(i)
	}
	//一直等待上面的线程跑完
	wg.Wait()
	//跑完了，继续跑后面的代码
	fmt.Println(c)
	//结果map[0:0 1:1 2:2 3:3 4:4 5:5 6:6 7:7 8:8 9:9]
}
func cover() {
	var a int = 89
	// 数据类型转换方式
	var b float32 = float32(a)

	fmt.Printf("%f ", b)

}
