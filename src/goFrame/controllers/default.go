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
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["index"] = &A{"username", 25}
	c.TplName = "index.tpl"
	
}

