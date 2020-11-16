package dal

import (
	"strconv"

	"ProductManage/public"
)

//产品的数据mode
type ProductModel struct {
	Id          string
	UserId      string
	CategoryId  string
	ProductName string
	Price       string
	Remark      string
	AddTime     string
}

//添加商品
func AddProduct(m *ProductModel) int {
	m.Id = public.ReplaceStr(m.Id)
	m.UserId = public.ReplaceStr(m.UserId)
	m.CategoryId = public.ReplaceStr(m.CategoryId)
	m.ProductName = public.ReplaceStr(m.ProductName)
	m.Price = public.ReplaceStr(m.Price)
	m.Remark = public.ReplaceStr(m.Remark)
	m.AddTime = public.ReplaceStr(m.AddTime)

	sql := "INSERT INTO  Product (id,UserId,CategoryId,ProductName,Price,Remark,AddTime)  VALUES "
	sql += "('" + public.GetUUIDS() + "','" + m.UserId + "','" + m.CategoryId + "','" + m.ProductName + "','" + m.Price + "','" + m.Remark + "','" + m.AddTime + "')"
	return ExecuteUpdate(sql)
}

//更新商品
func UpdateProduct(m *ProductModel) int {
	m.Id = public.ReplaceStr(m.Id)
	m.UserId = public.ReplaceStr(m.UserId)
	m.CategoryId = public.ReplaceStr(m.CategoryId)
	m.ProductName = public.ReplaceStr(m.ProductName)
	m.Price = public.ReplaceStr(m.Price)
	m.Remark = public.ReplaceStr(m.Remark)
	m.AddTime = public.ReplaceStr(m.AddTime)
	public.Log("UpdateProduct" + m.Id)
	sql := "UPDATE Product SET    CategoryId = '" + m.CategoryId + "' , ProductName = '" + m.ProductName + "' , Price = '" + m.Price + "' , Remark = '" + m.Remark + "'  where id= " + m.Id + ""
	return ExecuteUpdate(sql)
}

//删除
func DelProduct(id string) int {
	id = public.ReplaceStr(id)
	sql := "delete from  Product  where id= '" + id + "'"

	return ExecuteUpdate(sql)
}

//根据id获取商品
func GetProductIdByName(ProductName string) string {
	ProductName = public.ReplaceStr(ProductName)
	return GetSingle("SELECT id  FROM Product where ProductName='" + ProductName + "'")
}

//判断商品是否存在
func ExistsProductName(ProductName string) bool {
	ProductName = public.ReplaceStr(ProductName)
	ret := GetSingle("SELECT ProductName  FROM Product where  1=1   and  ProductName='" + ProductName + "'")
	if ret == "" {
		return false
	}
	return true
}

//获取最大的id
func GetProductMaxId() string {
	return GetSingle("SELECT max(id)  FROM Product ")
}

//获取商品总数
func GetProductTotal(wheresql string) string {
	sql := "SELECT count(id) as num  FROM  Product where  1=1   "
	if wheresql != "" {
		sql += wheresql
	}
	return GetSingle(sql)
}

//根据id获取商品json
func GetProductById(Id string) string {
	Id = public.ReplaceStr(Id)
	sql := "SELECT *  FROM Product a where Id= " + Id
	MapList := ExecuteQuery(sql)
	str := "{\"data\":[" + public.GetJsonStrByMap(MapList) + "]}"
	return str
}

//获取商品列表
func GetProductList(pageid string, pagesize string, userid string, ProductName string) string {
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
	ProductName = public.ReplaceStr(ProductName)
	sql := "SELECT *  FROM Product   where  1=1   "
	wheresql := ""
	if userid != "" {
		userid = public.ReplaceStr(userid)
		wheresql += " and userid = '" + userid + "'"
	}
	if ProductName != "" {
		ProductName = public.ReplaceStr(ProductName)
		wheresql += " and ProductName like '%" + ProductName + "%'"
	}
	sql += wheresql
	sql += "  order by id desc   limit " + pagesize + " OFFSET " + fromStr
	jsonString := ExecuteQueryJson(sql)
	total := GetProductTotal(wheresql)
	str := "{\"code\":0,\"total\":\"" + total + "\",\"curpage\":\"" + pageid + "\",\"pagesize\":\"" + pagesize + "\",\"data\":" + jsonString + "}"
	return str
}
