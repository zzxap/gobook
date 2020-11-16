package main

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}
func main() {

	go f("direct")
	//开启一个goroutine
	go f("goroutine")
	//开启一个goroutine
	go func(msg string) {
		fmt.Println(msg)
	}("going")
	//等待
	time.Sleep(time.Second)
	fmt.Println("done")
}

/*
打印结果  可以看到goroutine的执行是时间顺序是不确定的

D:\mybook\bookSource>go run goroutine.go
going
direct : 0
goroutine : 0
goroutine : 1
goroutine : 2
direct : 1
direct : 2
done

D:\mybook\bookSource>go run goroutine.go
direct : 0
direct : 1
direct : 2
goroutine : 0
goroutine : 1
goroutine : 2
going
done

D:\mybook\bookSource>go run goroutine.go
goroutine : 0
goroutine : 1
goroutine : 2
going
direct : 0
direct : 1
direct : 2
done

D:\mybook\bookSource>go run test.go
going
direct : 0
direct : 1
direct : 2
goroutine : 0
goroutine : 1
goroutine : 2
done

*/
