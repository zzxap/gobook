package main

import (
	"fmt"
	"time"
)

/*
interface 是GO语言的基础特性之一。可以理解为一种类型的规范或者约定。它跟java，C# 不太一样，
不需要显示说明实现了某个接口，它没有继承或子类或“implements”关键字，只是通过约定的形式，
隐式的实现interface 中的方法即可。因此，Golang 中的 interface 让编码更灵活、易扩展。

什么情况下使用interface ？
当我们给一个系统添加一个功能的时候，不是通过修改代码，而是通过增添代码来完成，那么就是开闭原则的核心思想了。
所以要想满足上面的要求，是一定需要interface来提供一层抽象的接口

作为interface数据类型，他存在的意义在哪呢？
实际上是为了满足一些面向对象的编程思想。目标就是高内聚，低耦合。

go中严格上说没有多态，但可以利用接口进行，
对于都实现了同一接口的两种对象，可以进行类似地向上转型，并且在此时可以对方法进行多态路由分发

*/
func main() {
	//测试
	testNilInterface()
	testInterface()

	/*
		输出:
		是自定义结构体类型
		student {abc}
		Woof!
		Meow!
		进行了奔跑业务...
		进行了睡觉业务...

	*/
}

//例子1 空接口
// 定义一个结构体
type Student struct {
	Name string
}

//空接口
func testNilInterface() {
	var v interface{}
	v = 12
	v = "abc"
	v = 12.22

	v = Student{Name: "abc"} // 自定义结构体类型
	//判断v的类型
	if _, ok := v.(int); ok {
		fmt.Printf("是int类型 \n")
	} else if _, ok := v.(string); ok {
		fmt.Printf(" 是字符串类型\n")
	} else if _, ok := v.(Student); ok {
		fmt.Printf(" 是自定义结构体类型\n")
	} else {
		fmt.Printf("未知类型\n")
	}
	//或者这样判断
	switch v.(type) {

	case bool:
		fmt.Printf("%s", v) //v
	case int:
		fmt.Printf("%d", v)
	case int64:
		fmt.Printf("%d", v)
	case int32:
		fmt.Printf("%d", v)
	case float64:
		fmt.Printf("%1.2f", v)
	case float32:
		fmt.Printf("%1.2f", v)
	case string:
		fmt.Printf("%s", v)
	case Student:
		fmt.Printf("student %s \n", v)
	case []byte:
		fmt.Printf("byte %s", string(v.([]byte)))
	case time.Time:
		fmt.Printf("%s", v)
	default:
		fmt.Printf("%s", v)
	}
}

/*
 interface 是一种具有一组方法的类型，这些方法定义了 interface 的行为
interface{} 会占用两个字长的存储空间，一个是自身的 methods 数据，一个是指向其存储值的指针，也就是 interface 变量存储的值
一个类型如果实现了一个 interface 的所有方法就说该类型实现了这个 interface，
空的 interface 没有方法，所以可以认为所有的类型都实现了 interface{}。
如果定义一个函数参数是 interface{} 类型，这个函数应该可以接受任何类型作为它的参数。
跟java c++ 其它语言的多态类似
*/

//例子2 非空接口
type Animal interface {
	Speak() string
}

type Dog struct {
}

//Dog实现了一个Animal 的Speak方法就说Dog类型实现了这个 Animal，
func (d Dog) Speak() string {
	fmt.Printf("Woof!\n")
	return "Woof!"
}

type Cat struct {
}

func (c *Cat) Speak() string {
	fmt.Printf("Meow!\n")
	return "Meow!"
}

func testInterface() {
	dog := Dog{}
	dog.Speak()
	cat := Cat{}
	cat.Speak()

	person := &Person{}
	person.Run()
	person.Sleep()

}

//例子3

//我们要写一个类Person
type Person struct {
}

//奔跑业务
func (p *Person) Run() {
	fmt.Println("进行了奔跑业务...")
}

//睡觉业务
func (p *Person) Sleep() {
	fmt.Println("进行了睡觉业务...")
}
