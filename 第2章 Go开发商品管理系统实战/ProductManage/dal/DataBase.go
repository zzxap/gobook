package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ProductManage/public"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB
var PthSep string
var ISmac int

//支持切换 pgsql sqlite
var Dbtype = "sqlite" // "pgsql" sqlite
var Dbhost = ""       //"127.0.0.1"
var Dbport = 5432
var Dbuser = "postgres"
var Dbpassword = ""
var Dbname = "invoice"

var oSingle sync.Once
var logsql = true

func GetTimeId() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func InitDB() bool {

	var testdb *sql.DB
	var err error
	if Dbtype == "pgsql" {

		public.Log("Init DB pgsql ...")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Dbhost, Dbport, Dbuser, Dbpassword, Dbname)
		testdb, err = sql.Open("postgres", psqlInfo)

	} else {
		public.Log("Init DB sqlite ...")
		curdir := public.GetCurDir()
		PthSep = string(os.PathSeparator)
		dbpath := curdir + PthSep + "db" + PthSep + "db.sqlite"

		if !public.ExistsPath(dbpath) {
			public.Log("db not exists" + dbpath)
		}
		if public.ExistsPath(dbpath) {
			public.Log(dbpath + "  存在")
			testdb, err = sql.Open("sqlite3", dbpath)
		} else {
			public.Log("db not exists" + dbpath)
		}
	}
	if err != nil {
		public.Log(err)
		return false
	}
	if testdb == nil {
		public.Log("Init DB fail")
		return false
	}
	public.Log("testing db connection...")
	err2 := testdb.Ping()
	public.Log("ping..." + Dbtype)
	if err2 != nil {
		fmt.Printf("Error on opening database connection: %s", err2.Error())
		return false
	} else {
		public.Log("connection.success")
		_db = testdb
		_db.SetMaxOpenConns(2000) //设置最大打开连接数
		_db.SetMaxIdleConns(10)   //设置最大空闲连接数
		return true
	}
}

//创建表
func AddTable() {
	sql := `CREATE TABLE IF NOT EXISTS  users (Id SERIAL PRIMARY KEY    NOT NULL , userId INTEGER, userName VARCHAR(50), password VARCHAR(50), userRole INTEGER, flag INTEGER, addTime timestamp(0) without time zone, DepartmentId INTEGER, phone  VARCHAR(50), address  VARCHAR(50), remark VARCHAR(50), ServerId INTEGER, dataFrom INTEGER, AddEditDel INTEGER DEFAULT 1, Power VARCHAR(50), Email VARCHAR(50), CardNum VARCHAR(50), IdNum VARCHAR(50), Sex VARCHAR(50), FilePath VARCHAR(50), CompanyName VARCHAR(50), CategoryId INTEGER, Name VARCHAR(50), Birthday VARCHAR(50));
CREATE TABLE IF NOT EXISTS  Product (Id SERIAL PRIMARY KEY    NOT NULL , name VARCHAR(50), phone VARCHAR(50), address VARCHAR(50), remarks VARCHAR(50), belong VARCHAR(50), ActualPay float, overdraft float, ServerId INTEGER, flag INTEGER, dataFrom INTEGER, userid INTEGER, AddEditDel INTEGER DEFAULT 1, Recharge FLOAT DEFAULT 0, Gifts FLOAT DEFAULT 0, CardNumber VARCHAR(50), Birthday timestamp(0) without time zone, Category INTEGER, Email VARCHAR(50), CardNum VARCHAR(50), IdNum VARCHAR(50), Sex VARCHAR(50), FilePath VARCHAR(50), Username VARCHAR(50));

`
	arr := strings.Split(sql, ";")
	for i := 0; i < len(arr); i++ {
		str := arr[i]
		str = strings.Replace(str, "\n", "", -1)
		if len(str) > 0 {
			ExecuteUpdateInDb(Getdb(), str)
		}
	}
}

func Getdb() *sql.DB {

	//err2 := _db.Ping()
	//if err2 != nil {
	//	log.Fatalf("Error on opening database connection: %s", err2.Error())
	//}
	// Ping验证连接到数据库是否还活着，必要时建立连接。

	return _db

}
func PrintCurrentPath() {

	dir, errer := filepath.Abs(filepath.Dir(os.Args[0]))
	if errer != nil {
		log.Fatal(errer)

	}
	public.Log(dir)
}

//获取 第一行第一列的json数据
func GetSingleJson(sqlStr string) string {

	return GetSingleJsonInDb(Getdb(), sqlStr)

}

//获取 第一行第一列的数据
func GetSingle(sqlStr string) string {

	return GetSingleInDb(Getdb(), sqlStr)

}

//获取 第一行第一列的数据
func GetSingleInDb(db *sql.DB, sqlStr string) string {
	var id string
	if logsql {
		public.Log(sqlStr)
	}

	rows, errr := db.Query(sqlStr)
	if errr != nil {
		public.Log(errr)
		return ""
	}
	//defer db.Close()  .Scan(&id)
	i := 0
	defer rows.Close()
	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		//public.Log("has row")
		i++
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
			public.Log("not Single result")
			return ""
		}
		//public.Log("333")
		for i, _ := range columns {
			v := values[i]
			if v == nil {
				id = ""
			} else {

				switch v.(type) {
				default:
					id = fmt.Sprintf("%s", v)
				case bool:
					id = fmt.Sprintf("%s", v) //v
				case int:
					id = fmt.Sprintf("%d", v)
				case int64:
					id = fmt.Sprintf("%d", v)
				case int32:
					id = fmt.Sprintf("%d", v)
				case float64:
					id = fmt.Sprintf("%1.2f", v)
				case float32:
					id = fmt.Sprintf("%1.2f", v)
				case string:
					id = fmt.Sprintf("%s", v)
				case []byte: // -- all cases go HERE!
					id = string(v.([]byte))
				case time.Time:
					id = fmt.Sprintf("%s", v)
				}

			}

		}
	}
	if i == 0 {
		//public.Log("has no row")
		return ""
	}

	//public.Log(id)
	return id
}

//获取 第一行第一列的数据
func GetSingleJsonInDb(db *sql.DB, sqlStr string) string {
	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return ""
	}
	defer rows.Close()

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

//执行sql 插入 更新 删除
func ExecuteUpdate(sqlStr string) int {

	return ExecuteUpdateInDb(Getdb(), sqlStr)

}

//执行sql 插入 更新 删除
func ExecuteUpdateInDb(db *sql.DB, sqlStr string) int {

	if strings.Contains(strings.ToLower(sqlStr), "insert") {
		sqlStr = strings.Replace(sqlStr, "'00:00:00'", "null", -1)

		if Dbtype == "pgsql" {
			rowId := 0
			sqlStr += " RETURNING id "

			err := db.QueryRow(sqlStr).Scan(&rowId)
			if err != nil {
				public.Log("exec sql failed:", err.Error()+" "+sqlStr)
				return 0
			}
			public.Log("lastInsertId=")
			public.Log(rowId)
			return rowId

		} else {
			res, err := db.Exec(sqlStr)
			if err != nil {
				public.Log("exec sql failed:", err.Error()+" "+sqlStr)
				return 0
			} else {
				//public.Log("exec Update sql success")
			}

			rowId, err := res.LastInsertId()
			if err != nil {
				public.Log("fetch last insert id failed:", err.Error())
				return 0
			}

			str := strconv.FormatInt(rowId, 10)
			ret, _ := strconv.Atoi(str)
			return ret
		}

	} else {
		res, err := db.Exec(sqlStr)
		if err != nil {
			public.Log("exec sql failed:", err.Error()+" "+sqlStr)
			return 0
		}
		rowId, err := res.RowsAffected()
		if err != nil {
			public.Log("fetch RowsAffected failed:", err.Error())
			return 0
		}
		str := strconv.FormatInt(rowId, 10)
		ret, _ := strconv.Atoi(str)
		return ret
	}
}

//查询返回map
func ExecuteQuery(sqlStr string) map[int]map[string]string {

	return ExecuteQueryInDb(Getdb(), sqlStr)
}

//查询数据 返回json字符串
func ExecuteQueryJson(sqlStr string) string {

	str := ExecuteQueryJsonInDb(Getdb(), sqlStr)
	if len(str) == 0 {
		return "[]"
	}
	return str
}
func GetRows(sqlStr string) *sql.Rows {
	db := Getdb()
	if db == nil {
		return nil
	}
	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return nil
	}

	return rows

}

//查询数据
func ExecuteQueryInDb(db *sql.DB, sqlStr string) map[int]map[string]string {

	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return nil
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)

	if count == 0 {
		return nil
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	final_result := make(map[int]map[string]string)
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
		final_result[result_id] = m
		result_id++
	}

	return final_result
}

//查询数据
func ExecuteQueryJsonInDb(db *sql.DB, sqlStr string) string {

	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return ""
	}
	defer rows.Close()

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
				case []byte:
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

func GetColumnName(sqlStr string) []string {
	return GetColumnNameInDb(Getdb(), sqlStr)
}

func GetColumnNameInDb(db *sql.DB, sqlStr string) []string {

	rows, err := db.Query(sqlStr)

	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		CheckErr(err)
		return nil
	}
	defer rows.Close()
	columns, _ := rows.Columns()

	names := make([]string, len(columns))
	for _, col := range columns {
		names = append(names, col)
	}
	return names
}

func GetColumnValueList(sqlStr string) []string {

	return GetColumnValueListInDb(Getdb(), sqlStr)

}
func GetColumnValueListInDb(db *sql.DB, sqlStr string) []string {

	rows, err := db.Query(sqlStr)

	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		CheckErr(err)
		return nil
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	names := make([]string, count)
	for rows.Next() {
		rows.Scan(valuePtrs...)

		var v interface{}
		val := values[0]
		b, err := val.([]byte)
		if err {
			v = string(b)
		} else {
			v = val
		}
		names = append(names, fmt.Sprintf("%s", v))
	}
	return names
}
func CheckErr(err error) {
	if err != nil {
		public.Log(err)
	}
}
