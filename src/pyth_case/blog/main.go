package main

import (
	_ "blog/routers"
	"github.com/astaxie/beego"
	"blog/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)



func init() {
	// 参数1   driverName
	// 参数2   数据库类型
	// 这个用来设置 driverName 对应的数据库类型
	// mysql / sqlite3 / postgres 这三种是默认已经注册过的，所以可以无需设置
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
	// 参数2        driverName
	// 参数3        对应的链接字符串
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
	maxIdle := 30
	maxConn := 30
	orm.RegisterDataBase("default", "mysql", "root:root@/tprbac?charset=utf8", maxIdle, maxConn)
}


func main() {
	// 开启 orm 调试模式：开发过程中建议打开，release时需要关闭
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("tprbac", false, true)

	// 创建一个 ormer 对象
	o := orm.NewOrm()
	o.Using("default")
	perfile := new(models.Profile)
	perfile.Age = 30

	user := new(models.User)
	user.Name = "tom"
	user.Profile = perfile

	// insert
	o.Insert(perfile)
	o.Insert(user)
	o.Insert(perfile)
	o.Insert(user)
	o.Insert(perfile)
	o.Insert(user)

	// update
	user.Name = "hezhixiong"
	num, err := o.Update(user)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)

	// delete
	o.Delete(&models.User{Id: 2})
	beego.Run()
}

