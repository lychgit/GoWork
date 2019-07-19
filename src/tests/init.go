/*
@Time : 2019/3/28 11:02
@Author : shilinqing
@File : Init
*/
package models

import (
	"github.com/shiLinQing407/BeeGoWeb/models/log"
	"github.com/shiLinQing407/BeeGoWeb/models/system"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)


var (
	db orm.Ormer
)

func Init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	timezone := beego.AppConfig.String("dbtimezone")
	if dbport == "" {
		dbport = "3306"
	}
	//fmt.Println(os.Getwd())
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	//註冊默認數據庫
	//fmt.Println(dsn)
	//orm.RegisterDataBase("default", "mysql", dsn)
	//orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/go_web?charset=utf8")
	//new用来分配内存，但与其他语言中的同名函数不同，它不会初始化内存，只会讲内存置零；
	// 也就是说，new(T)会为类型为T的新项分配已置零的内存空间，并返回他的地址，也就是一个类型为*T的值。
	// 用Go的术语来说，它返回一个指针，改指针指向新分配的，类型为T的零值
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/go?charset=utf8")
	orm.RegisterModelWithPrefix("go_", new(system.User), new(log.Log), new(system.Menu), new(system.Role), new(system.RoleMenuRel),new(system.RoleUserRel))
	//if beego.AppConfig.String("runmode") == "dev" {
	//	orm.Debug = true
	//}
	db = orm.NewOrm()
	orm.Debug = true
}

