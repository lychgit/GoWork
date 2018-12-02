package main

import (
	_ "goFrame/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/static", "static")	//设置静态文件处理目录
	//beego.SetViewsPath("templatePath") //设置模板目录
	beego.SetViewsPath("views") //设置模板目录
	beego.Run()
}

