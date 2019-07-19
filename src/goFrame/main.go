package main

import (
	//_ "goFrame/routers"
	"github.com/astaxie/beego"
	"goFrame/models"
	)

func init() {
	beego.Error("main.init")
	//beego.LoadAppConfig("ini", "conf/app.conf")
	models.Init()
	// 生产环境不输出debug日志
	if beego.AppConfig.String("runmode") == "dev" {
		//beego.SetLevel(beego.LevelInformational)
		beego.SetLevel(beego.LevelDebug)
	}
	beego.AppConfig.Set("version", beego.AppConfig.String("AppVer"))
}

func main() {
	beego.Run()
}
