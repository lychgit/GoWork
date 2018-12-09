package main

import (
	_ "goFrame/routers"
	"github.com/astaxie/beego"
	"strings"
)

func main() {
	beego.SetStaticPath("/static", "static")	//设置静态文件处理目录
	beego.SetStaticPath("/static/images","images")
	beego.SetStaticPath("/static/css","css")
	beego.SetStaticPath("/static/js","js")
	//beego.SetViewsPath("templatePath") //设置模板目录
	beego.SetViewsPath("views") //设置模板目录
	//当你设置了自动渲染，然后在你的 Controller 中没有设置任何的 TplName，那么 beego 会自动设置你的模板文件如下：
	//c.TplName = strings.ToLower(c.controllerName) + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
	//也就是你对应的 Controller 名字+请求方法名.模板后缀，也就是如果你的 Controller 名是 AddController，
	// 请求方法是 POST，默认的文件后缀是 tpl，那么就会默认请求 /viewpath/AddController/post.tpl 文件。
	beego.Run()
}

