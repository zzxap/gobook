/*
Sqlite 是一个本地数据库，无需安装，只需要一个文件
适合快速部署，对性能要求不是很高的场景

*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB
var PthSep string

var oSingle sync.Once
var logsql = true

func GetTimeId() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//初始化数据库
func InitDB() bool {

	var testdb *sql.DB
	var err error
	PthSep = string(os.PathSeparator)

	dbpath := "./db/db.sqlite"
	testdb, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Println(err)
		return false
	}
	if testdb == nil {
		log.Println("Init DB fail")
		return false
	}
	err2 := testdb.Ping()
	if err2 != nil {
		fmt.Printf("Error on opening database connection: %s", err2.Error())
		return false
	} else {
		_db = testdb
		_db.SetMaxOpenConns(2000) //设置最大打开连接数
		_db.SetMaxIdleConns(10)   //设置最大空闲连接数
		log.Println("Init DB success")
		return true
	}

}
func main() {

	if InitDB() {
		//拼接好sql插入，要过滤一些特殊字符防止sql注入
		//sql := "INSERT INTO  Users (id,UserName,Password,AddTime)  VALUES ('1','mm','123456','2020-06-06')"
		//ExecuteUpdate(_db, sql)
		//arr := []interface{"22", "mmmm22mm", "123456mmmm22", "2020-06-06"}
		arr := []interface{}{"13", "bb", "cc", "2020-06-06"}
		//用实务和参数化方式执行插入 可以防止sql注入
		sql := "INSERT INTO  Users (id,UserName,Password,AddTime)  VALUES (?,?,?,?)"
		ExecuteUpdateTran(_db, sql, arr)
		//查询
		jsonstr := ExecuteQuery(_db, "select * from users")
		log.Println(jsonstr)
		//结果 [{"addtime":"2020-06-06","id":"1","password":"123456","username":"mm"}]
	}

}

//…其实是go的一种语法糖。
//它的第一个用法主要是用于函数有多个不定参数的情况，可以接受多个不确定数量的参数。

//执行sql事务
func ExecuteUpdateTran(db *sql.DB, sqlStr string, argsList []interface{}) bool {
	//开启事务
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		fmt.Println("Prepare faillllll")
		return false
	}

	//将参数传递到sql语句中并且执行
	res, err := stmt.Exec(argsList...)
	if err != nil {
		fmt.Println("Exec fail2222")
		fmt.Println(err.Error())
		return false
	}
	//将事务提交
	tx.Commit()
	//获得上一个插入自增的id
	fmt.Println(res.LastInsertId())
	return true

}

//执行sql语句
func ExecuteUpdate(db *sql.DB, sqlStr string) int {

	res, err := db.Exec(sqlStr)
	if err != nil {
		log.Println("exec sql failed:", err.Error()+" "+sqlStr)
		return 0
	}
	rowId, err := res.RowsAffected()
	if err != nil {
		log.Println("fetch RowsAffected failed:", err.Error())
		return 0
	}

	str := strconv.FormatInt(rowId, 10)
	ret, _ := strconv.Atoi(str)
	return ret

}

//查询sql 返回json字符

func ExecuteQuery(db *sql.DB, sqlStr string) string {

	rows, err := db.Query(sqlStr)

	log.Println("sqlStr=" + sqlStr)

	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer rows.Close()
	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)

	if count == 0 {
		return ""
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	final_result := make([]map[string]string, 0)
	result_id := 0
	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		for i, col := range columns {
			v := values[i]
			key := strings.ToLower(col)
			if v == nil {
				m[key] = ""
			} else {

				switch v.(type) {
				default:
					m[key] = fmt.Sprintf("%s", v)
				case bool:

					m[key] = fmt.Sprintf("%s", v)
				case int:

					m[key] = fmt.Sprintf("%d", v)
				case int64:

					m[key] = fmt.Sprintf("%d", v)
				case float64:

					m[key] = fmt.Sprintf("%1.2f", v)
				case float32:

					m[key] = fmt.Sprintf("%1.2f", v)
				case string:

					m[key] = fmt.Sprintf("%s", v)
				case []byte: // -- all cases go HERE!

					m[key] = string(v.([]byte))
				case time.Time:
					m[key] = fmt.Sprintf("%s", v)
				}
			}
		}

		final_result = append(final_result, m)

		result_id++
	}

	jsonData, err := json.Marshal(final_result)
	if err != nil {
		return ""
	}

	return string(jsonData)

}
