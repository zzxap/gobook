package main

import (
	"fmt"
	"os"
)

func testRecover() {
	//全用defer + recover 来捕获和处理异常
	defer func() {
		err := recover() //recover内置函数可以捕获到异常
		if err != nil {  //nil是err的零值
			fmt.Println("err=", err)
			//runtime error: index out of range [3] with length 3
		}
	}() //匿名函数的调用方式一：func(){}()
	arr := []string{"a", "b", "c"}
	str := arr[3]
	fmt.Println("str=", str)
}
func main() {
	//测试
	test()

}

/*
Defer
Defer语句将一个函数放入一个列表（用栈表示其实更准确）中，该列表的函数在环绕defer的函数返回时会被执行。
defer通常用于简化函数的各种各样清理动作，例如关闭文件，解锁等等的释放资源的动作。

Panic
Panic是内建的停止控制流的函数。相当于其他编程语言的抛异常操作。当函数F调用了panic，F的执行会被停止，在F中panic前面定义的defer操作都会被执行，然后F函数返回。
对于调用者来说，调用F的行为就像调用panic（如果F函数内部没有把panic recover掉）。如果都没有捕获该panic，相当于一层层panic，程序将会crash。panic可以直接调用，也可以是程序运行时错误导致，例如数组越界。

Recover
Recover是一个从panic恢复的内建函数。Recover只有在defer的函数里面才能发挥真正的作用。如果是正常的情况（没有发生panic），
调用recover将会返回nil并且没有任何影响。如果当前的goroutine panic了，recover的调用将会捕获到panic的值，并且恢复正常执行。

Go语言追求简洁优雅，Go语言不支持传统的 try…catch…finally 这种异常，
Go语言的设计者们认为，将异常与控制结构混在一起会很容易使得代码变得混乱。

在Go语言中，使用多值返回来返回错误。不要用异常代替错误，更不要用来控制流程。
在极个别的情况下，也就是说，遇到真正的异常的情况下（比如除数为0了）。
才使用Go中引入的Exception处理：defer, panic, recover。

Go没有异常机制，但有panic/recover模式来处理错误
Panic可以在任何地方引发，但recover只有在defer调用的函数中有效

*/

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}

	written, err = io.Copy(dst, src)
	dst.Close()
	src.Close()
	//立马关闭
	return
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	//这里src不会立刻关闭，在函数return的时候关闭，用defer close防止忘记关闭
	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func testpanic() {
	//先声明defer，捕获panic异常
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获到了panic产生的异常：", err)
			fmt.Println("捕获到panic的异常了，recover恢复回来。")
		}
	}()
	//注意这个()就是调用该匿名函数的，
	//不写会报expression in defer must be function call
	/*
		 panic一般会导致程序挂掉（除非recover）  然后Go运行时会打印出调用栈
		但是，关键的一点是，即使函数执行的时候panic了，函数不往下走了，运行时并不是立刻向上传递panic，
		而是到defer那，等defer的东西都跑完了，panic再向上传递。所以这时候
		defer 有点类似 try-catch-finally 中的 finally。panic就是这么简单。
	*/
	panic("抛出一个异常了，defer会通过recover捕获这个异常，处理后续程序正常运行。")
	fmt.Println("这里不会执行了")
}
