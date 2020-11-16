package main

import (
	"context"
	"fmt"
	"time"
)

//此示例演示了如何使用可取消上下文来防止
// goroutine泄漏。 在示例函数结束时，goroutine开始了
// by gen将返回而不会泄漏。
func ExampleWithCancel() {
	// gen在单独的goroutine中生成整数，然后
	//将它们发送到返回的频道。
	// gen的调用者需要取消一次上下文
	//它们完成了对生成的整数的使用而不泄漏
	//内部goroutine由gen开始。
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return //返回不泄漏goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //当我们使用完整数后取消

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

//此示例传递具有任意截止日期的上下文以告知阻塞
//表示应该立即放弃工作的功能。
func ExampleWithDeadline() {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	//即使ctx将会过期，还是最好将其调用
	//在任何情况下都具有取消功能。 否则可能会使
	//上下文及其父对象的生存时间超出了必要。
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

	// Output:
	// context deadline exceeded
}

//此示例传递带有超时的上下文，以告知阻塞函数
//它应在超时后放弃工作。
func ExampleWithTimeout() {
	//传递带有超时的上下文，以告知阻塞函数
	//应该在超时结束后放弃工作。
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}

	// Output:
	// context deadline exceeded
}

//此示例演示如何将值传递到上下文
//以及如何检索它（如果存在）。
func ExampleWithValue() {
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, favContextKey("color"))

	// Output:
	// found value: Go
	// key not found: color
}
