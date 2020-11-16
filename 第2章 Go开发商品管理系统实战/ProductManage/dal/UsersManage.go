package dal

import (
	"strconv"

	"ProductManage/public"
)

type UsersModel struct {
	ID       string
	UserName string
	Password string
	AddTime  string
}

func AddUsers(m *UsersModel) int {
	m.ID = public.ReplaceStr(m.ID)
	m.UserName = public.ReplaceStr(m.UserName)
	m.Password = public.ReplaceStr(m.Password)
	m.AddTime = public.ReplaceStr(m.AddTime)

	sql := "INSERT INTO  Users (id,UserName,Password,AddTime)  VALUES "
	sql += "('" + public.GetUUIDS() + "','" + m.UserName + "','" + m.Password + "','" + m.AddTime + "')"
	return ExecuteUpdate(sql)
}
func UpdateUsers(m *UsersModel) int {
	m.ID = public.ReplaceStr(m.ID)
	m.UserName = public.ReplaceStr(m.UserName)
	m.Password = public.ReplaceStr(m.Password)
	m.AddTime = public.ReplaceStr(m.AddTime)

	sql := "UPDATE Users SET   UserName = '" + m.UserName + "' , Password = '" + m.Password + "' , AddTime = '" + m.AddTime + "'  where id= " + m.ID + ""
	return ExecuteUpdate(sql)
}
func DelUsers(id string) int {
	id = public.ReplaceStr(id)
	//sql := "delete from  Users  where id= '" + id + "'"
	sql := "update  Users set flag=0  where id= '" + id + "'"
	return ExecuteUpdate(sql)
}
func GetUsersIdByName(UserName string) string {
	UserName = public.ReplaceStr(UserName)
	return GetSingle("SELECT id  FROM Users where UserName='" + UserName + "'")
}
func ExistsUsersName(UserName string) bool {
	UserName = public.ReplaceStr(UserName)
	ret := GetSingle("SELECT UserName  FROM Users where  1=1   and  UserName='" + UserName + "'")
	if ret == "" {
		return false
	}
	return true
}
func GetUsersMaxId() string {
	return GetSingle("SELECT max(id)  FROM Users ")
}
func GetUsersTotal(wheresql string) string {
	sql := "SELECT count(id) as num  FROM  Users where  1=1   "
	if wheresql != "" {
		sql += wheresql
	}
	return GetSingle(sql)
}
func GetUsersById(Id string) string {
	Id = public.ReplaceStr(Id)
	sql := "SELECT *  FROM Users a where Id= " + Id
	MapList := ExecuteQuery(sql)
	str := "{\"data\":[" + public.GetJsonStrByMap(MapList) + "]}"
	return str
}
func GetUsersList(pageid string, pagesize string, UserName string) string {
	pageid = public.ReplaceStr(pageid)
	ipagesize, errr := strconv.Atoi(pagesize)
	if errr != nil {
		public.Log(errr)
		ipagesize = 50
		pagesize = "50"
	}
	p, _ := strconv.Atoi(pageid)
	var from int = (p - 1) * ipagesize
	fromStr := strconv.Itoa(from)
	UserName = public.ReplaceStr(UserName)
	sql := "SELECT *  FROM Users   where 1=1 "
	wheresql := ""

	if UserName != "" {
		UserName = public.ReplaceStr(UserName)
		wheresql += " and UserName like '%" + UserName + "%'"
	}
	sql += wheresql
	sql += "  order by id desc   limit " + pagesize + " OFFSET " + fromStr
	jsonString := ExecuteQueryJson(sql)
	total := GetUsersTotal(wheresql)
	str := "{\"code\":0,\"total\":\"" + total + "\",\"curpage\":\"" + pageid + "\",\"pagesize\":\"" + pagesize + "\",\"data\":" + jsonString + "}"
	return str
}
