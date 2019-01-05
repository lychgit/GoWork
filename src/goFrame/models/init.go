package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
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
	//dbprefix := beego.AppConfig.String("dbprefix")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	//註冊默認數據庫
	fmt.Println(dsn)
	orm.RegisterDataBase("default", "mysql", dsn)
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/go?charset=utf8")
	//不使用表名前缀
	//orm.RegisterModel(new(User), new(Task), new(TaskGroup), new(TaskLog))
	//使用表名前缀
	//dbprefix = "go_"
	//orm.RegisterModelWithPrefix(dbprefix, new(model.TaskLog))
	//orm.RegisterModelWithPrefix(dbprefix, new(User), new(Log))
	orm.RegisterModel(new(User), new(Log), new(Config), new(Menu), new(Role), new(RoleMenuRel),new(RoleUserRel))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/tprbac?charset=utf8")
	db = orm.NewOrm()
}

func TableName(name string) string {
	return beego.AppConfig.String("dbprefix") + name
}
