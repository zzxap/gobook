/*
第一步安装 protoc 在linux下操作

PROTOC_ZIP=protoc-3.7.1-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

下载一些必要的包
go get github.com/golang/protobuf/{proto,protoc-gen-go}

go get github.com/micro/protoc-gen-micro

新建文件  blueter.proto

protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. blueter.proto

如果提示

protoc-gen-micro: program not found or is not executable

执行
export GOROOT=/usr/local/go  #你的go安装路径
export GOPATH=$HOME/go       #go工程路径
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

执行
protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. blueter.proto
就会生成 blueter.pb.go  blueter.pb.micro.go  两个文件


把  blueter.pb.micro.go 中

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)
中的v2/ 去掉


go run server.go

如果提示panic: http: multiple registrations for /debug/requests

就是因为
github.com/coreos/etcd/vendor/golang.org/x/net/trace
和
golang.org/x/net/trace
冲突了

解决办法 删除掉vendor 里面的golang.org/x/net/trace

rm -rf $GOPATH/src/github.com/coreos/etcd/vendor/golang.org/x/net/trace



*/

package main

import (
	"log"
	"time"

	//proto "github.com/micro/examples/service/proto"
	"github.com/micro/go-micro"

	"./blueter"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"golang.org/x/net/context"
)

type Blueter struct{}

func (g *Blueter) Hello(ctx context.Context, req *blueter.HelloRequest, rsp *blueter.HelloResponse) error {
	rsp.Msg = "Hello ddddd " + req.From
	//这里可以做一些数据库新增修改查询操作把结果返回
	return nil
}

func main() {
	//使用ETCD作为注册中心
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"http://47.244.249.198:2379"}
	})

	service := micro.NewService(
		micro.Name("blueter"),
		micro.Version("laster"),
		micro.RegisterTTL(time.Second*30),      // 服务发现系统中的生存期
		micro.RegisterInterval(time.Second*10), ////重新注册的间隔

		//设置了30秒的TTL生存期，并设置了每10秒一次的重注册。

		micro.Registry(reg),
	)

	/*
		服务通过服务发现功能，在启动时进行服务注册，关闭时进行服务卸载。
		有时候这些服务可能会异常挂掉，进程可能会被杀掉，可能遇到短暂的网络问题。
		这种情况下，节点会在服务发现中被干掉。 理想状态是服务会被自动移除。

		解决方案
		为了解决这个问题，Micro注册机制支持通过TTL（Time-To-Live）和间隔时间注册两种方式。
		TTL指定一次注册在注册中心的有效期，过期后便删除，
		而间隔时间注册则是定时向注册中心重新注册以保证服务仍在线。

	*/

	service.Init()

	blueter.RegisterBlueterHandler(service.Server(), new(Blueter))
	//这样就自动注册了一个服务，如果在100台机器上运行了此程序，就注册了100个服务
	//客户端根据blueter请求就会自动负载均衡到这100台机器中的某一台
	//服务可以直接调用server.Run()来运行，这会让服务监听一个随机端口，这个调用也会让服务将自身注册到注册器，当服务停止运行时，会在注册器注销自己。
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

/*
运行
go run server.go
2020-06-05 10:43:14  file=auth/auth.go:31 level=info Auth [noop]
Authenticated as blueter-d37b2877-fe63-4744-b7a6-0501dcb2df42 issued by go.micro
2020-06-05 10:43:14  file=go-micro/service.go:206 level=info Starting [service] blueter
2020-06-05 10:43:14  file=grpc/grpc.go:864 level=info Server [grpc] Listening on [::]:60396
2020-06-05 10:43:14  file=grpc/grpc.go:697 level=info
 Registry [etcd] Registering node: blueter-d37b2877-fe63-4744-b7a6-0501dcb2df42

*/
