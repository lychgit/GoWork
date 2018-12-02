package routers

import (
	"blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/login",&controllers.LoginController{},"get:Login;post:Post")

	beego.AutoRouter(&controllers.AdminController{})
}
