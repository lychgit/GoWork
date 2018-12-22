package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/models"
	"goFrame/libs"
	"strings"
	"time"
	"strconv"
)

type LoginController struct {
	BaseController
}

// 登录
func (this *LoginController) Login() {
	if this.userId > 0 {
		beego.Debug("here is admin Login,uid=" + string(this.userId)) //debug埋点
		this.redirect("/admin")
	}
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()
		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if username != "" && password != "" {
			user, err := models.UserGetByName(username)
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now().Unix() //获取当前时间的Unix时间戳
				models.UserUpdate(user)
				beego.Debug("UserUpdate") //debug埋点
				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}
				this.redirect(beego.URLFor("AdminController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
		}
	}
	this.display()
}

// 退出登录
func (this *LoginController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor("AdminController.Login"))
}

// 获取系统时间
func (this *LoginController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}
