/*
MSSQL是微软数据库
优点是稳定可靠，缺点是收费
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

//初始化数据库
//Db数据库连接池
var _db *sql.DB
var PthSep string

var Dbhost = "127.0.0.1"
var Dbport = 1443
var Dbuser = "root"
var Dbpassword = "ap123"
var Dbname = "invoice"

func main() {

	if InitDB() {
		//插入操作
		//sql := "INSERT INTO  Users (id,UserName,Password,AddTime)  VALUES ('1','mm','123456','2020-06-06')"
		//ExecuteUpdate(_db, sql)
		//删除操作
		//ExecuteUpdate(_db, "delete from users where id=0")
		//查询
		jsonstr := ExecuteQuery(_db, "select * from mysql.user limit 1")
		log.Println(jsonstr)
		//数据库的insert, delete以及update的步骤内容大致一致，差别就是sql语句的变化
	}
}

func InitDB() bool {

	var testdb *sql.DB
	var err error
	log.Println("Init DB pgsql ...")
	//连接字符串
	connString := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", Dbhost, Dbport, Dbname, Dbuser, Dbpassword)

	testdb, err = sql.Open("mssql", connString)

	if err != nil {
		log.Println(err)
		return false
	}
	if testdb == nil {
		log.Println("Init DB fail")
		return false
	}

	log.Println("testing db connection...")

	err2 := testdb.Ping()

	if err2 != nil {
		fmt.Printf("Error on opening database connection: %s", err2.Error())
		return false
	} else {
		log.Println("connection.success")
		_db = testdb
		_db.SetMaxOpenConns(2000) //设置最大打开连接数
		_db.SetMaxIdleConns(10)   //设置最大空闲连接数

		return true
	}

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
