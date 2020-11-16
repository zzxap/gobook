package main

import "fmt"
import "time"

// 这是我们将要在多个并发实例中支持的任务了。这些执行者
// 将从 `jobs` 通道接收任务，并且通过 `results` 发送对应的结果。
// 将让每个任务间隔 1s 来模仿一个耗时的任务。
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

func main() {

	// 为了使用 worker 线程池并且收集他们的结果，我们需要
	// 2 个通道。
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// 这里启动了 3 个 worker，初始是阻塞的，因为
	// 还没有传递任务。
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 这里我们发送 9 个 `jobs`，然后 `close` 这些通道
	// 来表示这些就是所有的任务了。
	//如果不传递jobs <- j 任务会出现all goroutines are asleep
	//main goroutine线在等一个永远不会来的数据，那整个程序就永远等下去了

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// 最后，我们收集所有这些任务的返回值。
	for a := 1; a <= 9; a++ {
		<-results
	}
}

/*
执行这个程序，显示 9 个任务被多个 worker 执行。
整个程序 处理所有的任务仅执行了 3s 而不是 9s，是因为 3 个 worker 是并行的。

go run workerpools.go
输出
worker 3 processing job 1
worker 1 processing job 2
worker 2 processing job 3
worker 1 processing job 4
worker 2 processing job 5
worker 3 processing job 6
worker 1 processing job 7
worker 3 processing job 8
worker 2 processing job 9
*/
