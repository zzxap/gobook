package dal

//"ProductManage/public"

var sysdbname = "sqlite_master"
var columnname = "sql"
var tbname = "name"
var columnTable = "sqlite_master"
var autoinc = "AUTOINCREMENT"
var idtype = "INTEGER"
var datetime = "datetime"

func init() {
	if Dbtype == "pgsql" {
		sysdbname = "information_schema.tables"
		columnname = "column_name"
		tbname = "table_name"
		columnTable = "information_schema.columns"
		idtype = "SERIAL"
		autoinc = ""
		datetime = " timestamp(0) without time zone "
	}
}

func CreateTable() {

	sql := ""
	//pgsql如果不加双引号，那么表名都会被转化为小写。查询就不区分大小写
	num := GetSingle("SELECT count(*) as count  FROM " + sysdbname + "  where lower(" + tbname + ") ='product' ")
	if num == "0" {
		sql = "CREATE TABLE \"product\" (\"ID\" VARCHAR PRIMARY KEY NOT NULL ,\"UserId\" INTEGER, \"CategoryId\" INTEGER, \"ProductName\" VARCHAR, \"Price\"  FLOAT, \"AddTime\" " + datetime + ",  \"Remark\" VARCHAR)"
		ExecuteUpdate(sql)
	}
}
