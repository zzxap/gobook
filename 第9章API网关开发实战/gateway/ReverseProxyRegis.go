package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/creack/goproxy"
	"github.com/creack/goproxy/registry"
)

// ServiceRegistry is a local registry of services/versions
var ServiceRegistry = registry.DefaultRegistry{
	"service1": {
		"v1": {
			"localhost:9091",
			"localhost:9092",
		},
	},
}

func main() {

	http.HandleFunc("/", goproxy.NewMultipleHostReverseProxy(ServiceRegistry))
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%v\n", ServiceRegistry)
	})
	//增加服务  这里可以选择从一个数据库加载服务列表，再做一个管理系统管理这些服务即可
	ServiceRegistry.Add("service1", "v1", "10.10.10.1:9091")
	ServiceRegistry.Add("service2", "v1", "10.10.10.2:9091")
	ServiceRegistry.Add("service3", "v1", "10.10.10.3:9091")
	ServiceRegistry.Add("service1", "v1", "10.10.10.4:9091")
	//删除服务，可以用心跳做一个服务健康检查，服务是不可用或者超时，就可以删除掉

	ServiceRegistry.Delete("service1", "v1", "10.10.10.1:9091")
	ServiceRegistry.Delete("service2", "v1", "10.10.10.2:9091")
	ServiceRegistry.Delete("service3", "v1", "10.10.10.3:9091")
	ServiceRegistry.Delete("service1", "v1", "10.10.10.4:9091")
	//根据名称和版本号返回可用的服务列表
	ServiceRegistry.Lookup("service4", "v1")

	println("ready")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

/*
动态增加删除主机
// Registry is an interface used to lookup the target host
// for a given service name / version pair.
type Registry interface {
        Add(name, version, endpoint string)
          // Add an endpoint to our registry
        Delete(name, version, endpoint string)
         // Remove an endpoint to our registry
        Failure(name, version, endpoint string, err error)
         // Mark an endpoint as failed.
        Lookup(name, version string) ([]string, error)
         // Return the endpoint list for the given service name/version
}
*/
