/*
etcd是CoreOS团队于2013年6月发起的开源项目，
它的目标是构建一个高可用的分布式键值(key-value)数据库。
etcd内部采用raft协议作为一致性算法，etcd基于Go语言实现。

etcd作为服务发现系统，有以下的特点：

简单：安装配置简单，而且提供了HTTP API进行交互，使用也很简单
安全：支持SSL证书验证
快速：根据官方提供的benchmark数据，单实例支持每秒2k+读操作
可靠：采用raft算法，实现分布式系统数据的可用性和一致性


*/
package main

import (
	"RouteManage/public"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

func main() {
	InitDb()
}

var cli *clientv3.Client
var err error

func InitDb() {
	//log.Println("inidb")
	//etcd至少3台服务器，一台是主其它是从，数据一致，
	//如果其中一台崩了，其它两台会重新选举主从，保证服务可靠稳定
	etcdservers := "192.168.1.1:2379,192.168.1.2:2379,192.168.1.3:2379"
	array := strings.Split(etcdservers, ",")
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   array,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println("inidb fail")
		fmt.Println(err)
		return
	}

}

func Getdb() *clientv3.Client {
	if cli == nil {
		InitDb()
	}
	if cli != nil {
		//log.Println("client is not nil")
		return cli
	}
	log.Println("client is nil")
	return nil
}

//添加数据
func Put(key, value string) bool {
	//fmt.Println("put key=" + key + " value=" + value)
	if len(value) == 0 || len(key) == 0 {
		return false
	}

	Getdb()
	if cli == nil {
		fmt.Println("db init error")
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	resp, err := cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			fmt.Println("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			fmt.Println("client-side error: %v", err)
		default:
			fmt.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
		return false
	}
	if resp.PrevKv != nil {
		fmt.Println(resp.PrevKv)
	}

	//fmt.Println("success put key=" + key + " value=" + value)
	return true
}

//获取数据 根据key前缀获取相关的数据给返回一个map
func GetMap(key string) map[string]string {
	//fmt.Println("getmap key=" + key)

	Getdb()

	kv := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		log.Println("err %v", err)
		return kv
	}

	for _, ev := range resp.Kvs {

		kv[string(ev.Key)] = string(ev.Value)

	}
	return kv

}

//获取数据 根据key前缀获取相关的数据给返回一个map数组
func GetMapArray(key string) []map[string]string {
	Getdb()
	//fmt.Println("get key=" + key)
	//kv := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		log.Println("err %v", err)
		return nil
	}
	final_result := make([]map[string]string, 0)

	//log.Println(resp.Kvs)
	for _, ev := range resp.Kvs {
		m := make(map[string]string)
		m["key"] = string(ev.Key) // strings.Replace(string(ev.Key), key, "", -1)
		m["value"] = string(ev.Value)

		final_result = append(final_result, m)

	}
	return final_result

}

//删除数据，匹配key的前缀
func DeletePrefix(key string) bool {
	Getdb()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	_, err := cli.Delete(ctx, key, clientv3.WithPrefix()) //
	//withPrefix()是未了获取该key为前缀的所有key-value
	cancel()

	if err != nil {
		return false
	}

	return true

}

//删除数据匹配整个key
func Delete(key string) bool {
	Getdb()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	_, err := cli.Delete(ctx, key) //, clientv3.WithPrefix()
	//withPrefix()是未了获取该key为前缀的所有key-value
	cancel()

	if err != nil {
		return false
	}

	return true

}

//监控key的数据变化
func Watch(key string) {
	wc := cli.Watch(context.Background(), key, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wc {
		if v.Err() != nil {
			//panic(err)
		}
		for _, e := range v.Events {
			fmt.Printf("type:%v\n kv:%v  prevKey:%v  ", e.Type, e.Kv, e.PrevKv)
		}
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
