/*
golang中的锁是通过CAS原子操作实现的，Mutex结构如下：
type Mutex struct {
    state int32
    sema  uint32
}

//state表示锁当前状态，每个位都有意义，零值表示未上锁
//sema用做信号量，通过PV操作从等待队列中阻塞/唤醒goroutine，
等待锁的goroutine会挂到等待队列中，并且陷入睡眠不被调度，unlock锁时才唤醒。
具体在sync/mutex.go Lock函数实现中。

使用原子操作可以跨多个协程管理简单的计数器状态。
对于更复杂的状态，我们可以使用互斥锁 安全地跨多个协程访问数据。
sync.Mutex不区分读写锁，只有Lock()与Lock()之间才会导致阻塞的情况，
如果在一个地方Lock()，在另一个地方不Lock()而是直接修改或访问共享数据，
这对于sync.Mutex类型来说是允许的，因为mutex不会和goroutine进行关联。
如果想要区分读、写锁，可以使用sync.RWMutex类型

在Lock()和Unlock()之间的代码段称为资源的临界区(critical section)，
在这一区间内的代码是严格被Lock()保护的，是线程安全的，
任何一个时间点都只能有一个goroutine执行这段区间的代码。

尽量减少锁的持有时间，毕竟使用锁是有代价的，通过减少锁的持有时间来减轻这个代价：
细化锁的粒度。通过细化锁的粒度来减少锁的持有时间以及避免在持有锁操作的时候做各种耗时的操作。
不要在持有锁的时候做 IO 操作。尽量只通过持有锁来保护 IO 操作需要的资源而不是 IO 操作本身

*/

package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	//在我们的示例中，“状态”将是map。
	var state = make(map[int]int)
	//这个`mutex`将同步对`state`的访问。
	var mutex = &sync.Mutex{}
	//跟踪多少读写操作
	var readOps uint64
	var writeOps uint64
	//开启100个协程以执行重复读取状态，每毫秒一次
	for r := 0; r < 100; r++ {
		go func() {
			total := 0
			for {

				key := rand.Intn(5)
				//以独占方式访问状态 在这一区间内的代码是严格被Lock()保护的，是线程安全的
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)
				//等待两次读取之间的时间。
				//time.Sleep(time.Millisecond)
				runtime.Gosched()
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)
				//time.Sleep(time.Millisecond)
				// 为了确保这个 Go 协程不会在调度中饿死，我们
				// 在每次操作后明确的使用 `runtime.Gosched()`
				// 进行释放。这个释放一般是自动处理的

				//runtime.Gosched()用于让出CPU时间片。
				//这就像跑接力赛，A跑了一会碰到代码runtime.Gosched()就把接力棒交给B了

				runtime.Gosched()
			}
		}()
	}

	//让10个协程在`state`上操作一秒钟，
	time.Sleep(time.Second)
	//获取并报告最终操作计数。
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
	//最终锁定状态为`state`，说明它如何结束。
	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}

/*
运行结果
readOps: 83285
writeOps: 8320
state: map[1:97 4:53 0:33 2:15 3:2]
*/
