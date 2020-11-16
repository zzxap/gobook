package main

import (
	"context"
	"fmt"

	"./blueter"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func main() {
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"http://47.244.249.198:2379",
		}
	})

	service := micro.NewService(
		micro.Registry(reg),
	)
	service.Init()
	//blueter 是服务端注册的服务
	client := blueter.NewBlueterService("blueter", service.Client())
	param := &blueter.HelloRequest{
		From: "client",
		To:   "server",
		Msg:  "hello xxx",
	}

	rsp, err := client.Hello(context.Background(), param)
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp)
}
