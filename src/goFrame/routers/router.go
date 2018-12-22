package routers

import (
	admin "goFrame/admin/controllers"
	web "goFrame/web/controllers"
	"goFrame/controllers"
	"github.com/astaxie/beego"
	"net/http"
	"html/template"
)

func init() {
	// 设置默认404页面
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(rw, data)
	})
	/*路由设置*/

	//固定路由
	beego.Router("/", &controllers.MainController{}, "*:Index")
	//beego.Router("/", &controll.MainController{}, "*:Index")
	//beego.Router("/profile", &controllers.LoginController{}, "*:Profile")
	//beego.Router("/gettime", &controllers.LoginController{}, "*:GetTime")
	beego.Router("/login", &web.LoginController{}, "*:Login")
	beego.Router("/logout", &web.LoginController{}, "*:Logout")
	//beego.Router("/help", &controllers.HelpController{}, "*:Index")


	beego.Router("/admin", &admin.AdminController{}, "*:Index")
	beego.Router("/admin/login", &admin.LoginController{}, "*:Login")
	beego.Router("/admin/logout", &admin.LoginController{}, "*:Logout")
	//beego.Router("/admin/profile", &controllers.AdminController{}, "*:Profile")
	//beego.Router("/admin/gettime", &controllers.AdminController{}, "*:GetTime")

	//路由自动匹配
	/******************************		後台控制器		******************************/
	//beego.AutoRouter(&admincontroller.AdminController{})

	//beego.Router("/craw_movie", &controllers.CrawMovieController{},"*:CrawMovie")

}
