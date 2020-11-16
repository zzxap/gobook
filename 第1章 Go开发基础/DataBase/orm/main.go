package main

import (
	//"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//初始化数据库
//Db数据库连接池
var _db *gorm.DB
var PthSep string

func main() {

	if InitDB() {
		testdata(_db)
	}
}

//表结构
type UserInfo struct {
	Id       int    `gorm: "primary_key;AUTO_INCREMENT:number"`
	UserName string `gorm:"column:username"`
	Password string
}

func InitDB() bool {

	var err error
	log.Println("Init DB  ...")
	//构建连接, 格式是："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"

	_db, err = gorm.Open("mysql", "root:Sky_1728@tcp(192.168.1.74:3306)/mysql?charset=utf8")
	if err != nil {
		log.Println(err)
		return false
	}
	if _db == nil {
		log.Println("Init DB fail")
		return false
	}

	log.Println("connection.success")
	// 全局禁用表名复数
	_db.SingularTable(true)

	has := _db.HasTable(&UserInfo{})
	if !has {
		_db.AutoMigrate(&UserInfo{})
		fmt.Println("创建表")
	}

	return true

}

func testdata(db *gorm.DB) {
	//添加数据
	db.Create(&UserInfo{UserName: "lulu", Password: "123"})
	db.Create(&UserInfo{UserName: "cc", Password: "123"})

	//查询数据

	//获取第一条记录
	fmt.Println("获取第一条记录")
	var user UserInfo
	db.First(&user)
	fmt.Println(user)

	//更新数据
	fmt.Println("更新lulu全部字段")
	user.UserName = "lulu"
	user.Password = "111111"
	db.Save(user)

	fmt.Println("更新部分字段")
	db.Model(&user).Update("username", "luuu")

	//删除记录
	var del_user UserInfo
	del_user.Id = 4
	del_user.UserName = "cc"
	db.Delete(&del_user)

	//获取所有记录
	fmt.Println("获取所有记录")
	var users []UserInfo
	db.Find(&users)
	fmt.Println(users)

	//查询
	var userinfo UserInfo
	err := _db.Debug().Model(&userinfo).Where("username = ?", "luuu").Scan(&userinfo).Error
	if err != nil {
		return
	}
	fmt.Println(userinfo)

}
