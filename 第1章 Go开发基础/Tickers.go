//timer和ticker都是用于计时的
//使用timer定时器，超时后需要重置，才能继续触发。
//ticker只要定义完成，从此刻开始计时，不需要任何其他的操作，每隔固定时间都会触发。

package main

import (
    "fmt"
    "time"
)

func main() {

    timer1 := time.NewTimer(2 * time.Second)

    <-timer1.C
    fmt.Println("Timer 1 fired")

    timer2 := time.NewTimer(time.Second)
    go func() {
        <-timer2.C
        fmt.Println("Timer 2 fired")
    }()
    stop2 := timer2.Stop()
    if stop2 {
        fmt.Println("Timer 2 stopped")
    }

    time.Sleep(2 * time.Second)
}



package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}


//Go 中 最重要的状态管理方式是通过通道间的沟通来完成的，在 worker-pools中碰到过，但是
//还是有一些其他的方法来管理状态。使用`sync/atomic`包在多个 Go 协程中进行_原子计数_
package main

import (
    "fmt"
    "runtime"
    "sync/atomic"
    "time"
)

func main() {

    //使用一个无符号整型数来表示（永远是正整数）这个计数器
    var ops uint64 = 0

    //为了模拟并发更新，启动50个Go协程，对计数器每隔1ms进行一次加一操作
    for i := 0; i < 50; i++ {
        go func() {
            for {
                // 使用`AddUint64`来让计数器自动增加，使用`&`语法来给出`ops`的内存地址
                atomic.AddUint64(&ops, 1)

                //允许其他 GO 协程的执行
                runtime.Gosched()
            }
        }()
    }

    //等待一秒，让ops的自加操作执行一会
    time.Sleep(time.Second)

    //为了让计数器还在被其它  GO 协程更新时， 安全的使用它，
    //通过`LoadUint64`将当前的值拷贝提取到`opsFinal`中。
    opsFinal := atomic.LoadUint64(&ops)
    fmt.Println("ops:", opsFinal)
}

