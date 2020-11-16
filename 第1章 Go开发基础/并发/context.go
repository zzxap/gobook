/*
在golang中的创建一个新的协程并不会返回像c语言创建一个线程一样类似的pid，
这样就导致我们不能从外部杀死某个线程，所以我们就得让它自己结束。
（备注：goroutine不能返回pid的原因，应该是协程的实现原理有很大关系，多个协程对应1个线程的实现机制。）

当然我们可以采用channel＋select的方式，来解决这个问题，不过场景很复杂的时候，
我们就需要花费很大的精力去维护channel与这些协程之间的关系，这就导致了我们的并发代码变得很难维护和管理。
例如：由一个请求衍生出多个协程，并且之间需要满足一定的约束关系，以实现一些诸如：有效期，中止线程树，
传递请求全局变量之类的功能。

Context机制：context的产生，正是因为协程的管理问题，golang官方从1.7之后引入了context，
用来专门管理协程之间的关系。

Google的解决方法是Context机制，相互调用的goroutine之间通过传递context变量保持关联，
这样在不用暴露各goroutine内部实现细节的前提下，有效地控制各goroutine的运行。通过传递context就可以追踪goroutine
调用树，并在这些调用树之间传递通知和元数据。


虽然goroutine之间是平行的，没有继承关系，但是Context设计成是包含父子关系的，
这样可以更好的描述goroutine调用之间的树型关系。


context包的核心就是Context接口，其定义如下：

type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
Deadline方法    返回一个超时时间。到了该超时时间，该Context所代表的工作将被取消继续执行。Goroutine获得了超时时间后，可以对某些io操作设定超时时间。

Done方法    返回一个通道（channel）。当Context被撤销或过期时，该通道被关闭。它是一个表示Context是否已关闭的信号。

Err方法    当Done通道关闭后，Err方法返回值为Context被撤的原因。

Value方法    可以让Goroutine共享一些数据，当然获得数据是协程安全的。但使用这些数据的时候要注意同步，比如返回了一个map，而这个map的读写则要加锁。

注意：context包里的方法是线程安全的，可以被多个线程使用。

Context接口没有提供方法来设置其值和过期时间，也没有提供方法直接将其自身撤销。也就是说，Context不能改变和撤销其自身

*/

package main

import (
	"fmt"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	fmt.Println("server: hello handler started")
	defer fmt.Println("server: hello handler ended")

	select {
	case <-time.After(10 * time.Second):
		fmt.Fprintf(w, "hello\n")
	case <-ctx.Done():

		err := ctx.Err()
		fmt.Println("server:", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8090", nil)
}

func httpContext() {
	// 创建一个http服务倾听端口8000
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Fprint(os.Stdout, "processing request\n")
		// We use `select` to execute a peice of code depending on which
		// channel receives a message first
		select {
		case <-time.After(2 * time.Second):
			// If we receive a message after 2 seconds
			// that means the request has been processed
			// We then write this as the response
			w.Write([]byte("请求处理"))
		case <-ctx.Done():
			// If the request gets cancelled, log it
			// to STDERR
			fmt.Fprint(os.Stderr, "请求已取消\n")
		}
	}))
}

/*
您可以通过运行服务器并在浏览器中打开localhost：8000进行测试。
如果您在2秒钟前关闭浏览器，则应该在终端窗口上看到“请求已取消”字样。
*/

//context控制协程的例子
package main

import (
    "sync"
    "fmt"
    "time"
    "context"
    "math/rand"
)

func test1(wg * sync.WaitGroup, ctx context.Context, cfunc context.CancelFunc){
    defer wg.Done()
    fmt.Println("test1 run")
    d := time.Duration(rand.Intn(10))
    fmt.Printf("%d 秒之后，test1 将终止\n", d)
    time.AfterFunc( d * time.Second, cfunc)
    for{
        select {
        case <-ctx.Done():
            fmt.Println("test1 return")
            return
        default:
            fmt.Println("default")
            time.Sleep(1 * time.Second)

        }
    }
}

func test2(wg *sync.WaitGroup, ctx context.Context, cfunc context.CancelFunc){
    defer wg.Done()
    fmt.Println("test2 run")
    d := time.Duration(rand.Intn(10))
    fmt.Printf("%d 秒之后，test2 将终止\n", d)
    time.AfterFunc( d * time.Second, cfunc)
    for {
        select {
        case <-ctx.Done():
            fmt.Println("test2 return")
            // 这里要不用return 要不就用break + lebal， 不能直接用break，只用break，这里只跳出case
            return
        default:
            fmt.Println("test2 default")
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    wg := sync.WaitGroup{}
    rand.Seed(time.Now().UnixNano())
    ctx, cancel := context.WithCancel(context.Background())
    wg.Add(2)
    timeStart := time.Now().Unix()
    go test1(&wg, ctx, cancel)
    go test2(&wg, ctx, cancel)
    wg.Wait()
    timeEnd := time.Now().Unix()
    fmt.Printf("运行时长为：%d s\n", timeEnd - timeStart)
    fmt.Println("主协成退出!")
}