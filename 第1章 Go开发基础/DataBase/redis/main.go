/*
Redis是一个开源的、使用C语言编写的、支持网络交互的、可基于内存也可持久化的Key-Value数据库。

Redis 优势
性能极高 – Redis能读的速度是110000次/s,写的速度是81000次/s 。

丰富的数据类型 – Redis支持二进制案例的 Strings, Lists, Hashes, Sets 及 Ordered Sets 数据类型操作。

原子 – Redis的所有操作都是原子性的，同时Redis还支持对几个操作全并后的原子性执行。

丰富的特性 – Redis还支持 publish/subscribe, 通知, key 过期等等特性。


*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	//"github.com/go-redis/redis"
)

var conn redis.Conn
var err error

func main() {
	fmt.Println("redis set")
	Set("key11", "value11")
	Get("key11")
	fmt.Println("redis set finish")
}
func init() {
	fmt.Println("redis init")
	conn, err = redis.Dial("tcp", "192.168.1.74:6379")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//defer conn.Close()
}

func Set(key, value string) bool {
	if conn == nil {
		fmt.Println("redis nil")
		return false
	}
	// _, err = conn.Do("SET", "key", "value", "EX", "5") 设置过期 秒,过期后再读就是nil了
	_, err = conn.Do("SET", key, value)
	if err != nil {
		fmt.Println("redis set failed:", err.Error())
		return false
	}
	return true
}

func Get(key string) string {
	if conn == nil {
		fmt.Println("redis nil")
		return ""
	}
	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return ""
	} else {
		fmt.Printf("Get key= %v \n", value)
		return value
	}
}

func ExistsKey(key string) bool {
	if conn == nil {
		fmt.Println("redis nil")
		return false
	}
	is_key_exit, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		fmt.Println("error:", err)
		return false
	} else {
		fmt.Printf("exists or not: %v \n", is_key_exit)
		return true
	}

}
func DeleteKey(key string) bool {
	if conn == nil {
		fmt.Println("redis nil")
		return false
	}
	_, err = conn.Do("DEL", "mykey")
	if err != nil {
		fmt.Println("redis delelte failed:", err)
		return false
	}
	return true
}

func SetJson(key, value string) bool {

	if conn == nil {
		fmt.Println("redis nil")
		return false
	}
	n, err := conn.Do("SETNX", key, value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if n == int64(1) {
		fmt.Println("success")
	}
	return true
}

func GetJson(key string) string {
	if conn == nil {
		fmt.Println("redis nil")
		return ""
	}
	var imap map[string]string

	valueGet, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		fmt.Println(err)
	}

	errShal := json.Unmarshal(valueGet, &imap)
	if errShal != nil {
		fmt.Println(err)
	}
	fmt.Println(imap["username"])
	fmt.Println(imap["password"])

	return string(valueGet[:])
}
