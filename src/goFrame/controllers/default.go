package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type A struct {
	Name string
	Age int
}
func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.Data["index"] = &A{"username", 25}
	this.Layout = "index_layout.tpl"
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["HtmlHead"] = "blogs/html_head.tpl"
	//this.LayoutSections["Scripts"] = "blogs/scripts.tpl"
	//this.LayoutSections["Sidebar"] = ""
	this.TplName = "index.tpl"
}

