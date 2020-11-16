package mypack

import (
	"fmt"
)

//公共函数首字母大写
//在其它package内通过 包名.公有函数名 调用 比如 mypack.PublicFunc()
func PublicFunc() {

}

//私有函数首字母小写  在其它包内无法通过 包名.函数名 调用
func privateFunc() {

}

//公共变量  在同一个包内可以直接访问调用
//在其它包内可以通过 包名.公共变量 调用
var PublicString string

//私有变量  在同一个包内可以直接访问调用
var privateString string

/*
普通函数声明（定义）
函数声明包括函数名、形式参数列表、返回值列表（可省略）以及函数体。

func 函数名(形式参数列表)(返回值列表){
    函数体
}

形式参数列表描述了函数的参数名以及参数类型，这些参数作为局部变量，
其值由参数调用者提供，返回值列表描述了函数返回值的变量名以及类型，如果函数返回一个无名变量或者没有返回值，返回值列表的括号是可以省略的。
*/
//只有一个返回值
func add(a, b, c int) int {
	return a + b + c
}

//多个返回值
func cost(a, b int, name string) (int, string) {
	if name == "li" {
		return a + b + 100, "hight"
	} else {
		return a + b + 10, "low"
	}
}

//调用
result,title:=cost(2,3,"li")

//无返回值的函数
func cost(a,b int){
	
}


Go 函数方法
      在 Go 语言中，函数和方法不太一样，有明确的概念区分。其他语言中，比如 PHP 函数就是方法，方法

就是函数，但在 Go 语言中，函数是不属于任何结构体、类型的方法，也就是说函数是没有接收者的；而方法

是有接收者的，我们所说的方法要么属于一个结构体的，要么属于一个新定义的类型的

函数

       函数和方法，虽然概念不同，但是定义非常相似。函数的定义声明没有接收者，如下示例：

1
2
3
4
5
6
7
8
9
10
11
12
13
package main
 
import "fmt"
 
// 定义一个求两数之和的函数
func add(a,b int) int  {
    return a + b
}
 
func main() {
    sum := add(1,2)
    fmt.Println(sum)
}
　　上面的例子中，我们定义了一个 add 函数，它的函数签名是 func add(a,b int) int ，没有接收者，直接定义在 go 

的一个包之下，可以直接用，例子中的这个函数名是小写的 add，所以它的作用域只有当前包，不能被其他包使用，

如果我们把函数名以大写开头，可以被其他包调用，这也是 Go 中大小写的作用

方法 
       方法的声明和函数类似，她们的区别是：方法在定义的时候在 func 和 方法名之间有一个参数，这个参数就是方

法的接收者，这样我们定义的这个方法就和接收者绑定在一起了吧，称之为这个接收者的方法：

1
2
3
4
5
6
7
8
// 定义一个结构体
type person struct {
    name string
}
 
func (p person) String() string  {
    return "the person name is " + p.name
}
　　上面的这个例子中，func 和 方法名 String 之间的参数(p person) 就是接收者，现在我们说，类型 person 有了一个方法

String，现在看下如何使用：

func main()  {
    p := person{name : "张三"}
 
    fmt.Println(p.String())
}
　　Go 语言中接收者分为两种类型，值接受者和指针接收者。我们上面的例子中，就是值类型接收者的示例

       使用值类型接收者定义的方法，在调用的时候，使用的其实就是值接受者的一个副本，所以对该值的任何操作，不会影响

原来的值

package main
 
import "fmt"
 
type person struct {
    name string
}
 
func (p person) String() string {
    return "the person name is " + p.name
}
 
func (p person) modify()  {
    p.name = "李四"
}
 
func main()  {
    p := person{"张三"}
    // 值类型接收者
    p.modify()
 
    fmt.Println(p.String())
}
　　上面的例子中，打印出来的值是 "张三"，对其进行的修改无效。如果我们使用一个指针作为接收者，那么就会起作用了，

因为指针接收者传递的是一个指向原值指针的拷贝，指针的副本，指向的还是原来类型的值，所以修改时，同时也会影响原

来类型变量的值

package main
 
import "fmt"
 
type person struct {
    name string
}
 
func (p person) String() string {
    return "the person name is " + p.name
}
 
func (p *person) modify()  {
    p.name = "李四"
}
 
func main()  {
    p := person{"张三"}
    // 值类型接收者
    p.modify()
 
    fmt.Println(p.String())
}
　　  在调用方法的时候，传递的接收者本质上都是副本，只不过一个是值副本，另一个是指向这个值指针的副本。指针具有

指向原值的特性，所以修改了指针指向的值，也就修改了原有的值。我们可以简单的理解为值接收者使用的是值的副本来调用

方法，而指针接收者使用的是实际值调用方法

         上面的例子中我们发现，在调用指针接收者方法的时候，使用的也是一个值的变量，并不是一个指针，修改如下：

1
2
p := person{"张三"}
(&p).modify()
　　这样也是可以的，如果我们没有强制使用指针进行调用，Go 编译器会帮我们取指针，同样的，如果是一个值接收者的方法

使用指针也可以调用，Go 编译器会自动解引，如下：

1
2
p := person{"张三"}
fmt.Println((&p).String())
　　所以，方法的调用既可以是值也可以是指针

多值返回

      Go 语言支持函数方法的多值返回，也就是说我们定义的函数方法支持可以返回多个值，比如标准库里的很多方法，都是返回两

个值，第一个是函数需要返回的值，第二个是出错时返回的错误信息


package main
 
import (
    "fmt"
    "log"
    "os"
)
 
func main() {
    file , err := os.Open("/usr/tmp")
 
    if err != nil {
        log.Fatal(err)
        return
    }
 
    fmt.Println(file)
}
　　如果返回的值，我们不想使用，可以使用 _ 进行忽略：

1
file , _ := os.Open("/usr/tmp/")
　　多个值返回的定义也非常简单，如下示例：

func add(a,b int) (int,error) {
     return a + b, nil
}
　　函数方法声明定义的时候，采用逗号分隔，因为是多个返回，还要用括号扩起来，返回的值还是使用 return 关键字，以逗号分隔，和

声明的顺序一致

可变参数

       函数方法的参数可以是任意多个，这种我们称之为可变参数，比如我们常用的 fmt.Println() 这类函数，可以接收可变参数


func main() {
    fmt.Println("1","2","3","4")
}
自己定义一个可接收可变参数的函数，如果可变参数的类型是一样的则可以使用省略号 ... 代替：


func print(a ...interface{})  {
    for _,v := range a{
        fmt.Println(v)
    }
}
　　可变参数本质上是一个数组，所以我们可以像数组一样使用它


func reflect_typeof(a interface{}) {
    t := reflect.TypeOf(a)
    fmt.Printf("type of a is:%v\n", t)
 
    k := t.Kind()
    switch k {
    case reflect.Int64:
        fmt.Printf("a is int64\n")
    case reflect.String:
        fmt.Printf("a is string\n")
    }
}

func main() {
    var x int64 = 3
    reflect_example(x)
 
    var y string = "hello"
    reflect_example(y)
}
打印结果：

type of a is:int64
a is int64
type of a is:string
a is string

func reflect_value(a interface{}) {
    v := reflect.ValueOf(a)
    k := v.Kind()
    switch k {
    case reflect.Int64:
        fmt.Printf("a is Int64, store value is:%d\n", v.Int())
    case reflect.String:
        fmt.Printf("a is String, store value is:%s\n", v.String())
    }
}



Field()利用反射获取结构体里面的方法和调用。

 获取结构体的字段

我们可以通过上面的方法判断一个变量是不是结构体。

可以通过 NumField() 获取所有结构体字段的数目、进而遍历，通过Field()方法获取字段的信息。


type Student struct {
    Name  string
    Sex   int
    Age   int
    Score float32
}
 
func main() {
    //创建一个结构体变量
    var s Student = Student{
        Name:  "orange",
        Sex:   1,
        Age:   10,
        Score: 80.1,
    }
 
    v := reflect.ValueOf(s)
    t := v.Type()
    kind := t.Kind()
     
    //分析s变量的类型，如果是结构体类型，那么遍历所有的字段
    switch kind {
 
    case reflect.Struct:
        fmt.Printf("s is struct\n")
        fmt.Printf("field num of s is %d\n", v.NumField())
        //NumFiled()获取字段数，v.Field(i)可以取得下标位置的字段信息，返回的是一个Value类型的值
        for i := 0; i < v.NumField(); i++ {
            field := v.Field(i)
            //打印字段的名称、类型以及值
            fmt.Printf("name:%s type:%v value:%v\n",
                t.Field(i).Name, field.Type().Kind(), field.Interface())
        }
    default:
        fmt.Printf("default\n")
    }
}
执行结果：

s is struct
field num of s is 4
name:Name type:string value:orange
name:Sex type:int value:1
name:Age type:int value:10
name:Score type:float32 value:80.1