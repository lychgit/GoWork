package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/admin/models"
	BaseModels "goFrame/models"
	"strings"
	"strconv"
	"goFrame/libs"
	"time"
	"fmt"
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	user           *models.User
	userId         int
	userName       string
	pageSize       int
}

//这个函数主要是为了用户扩展用的，这个函数会在Get、Post、Delete、Put、Finish等这些 Method 方法之前执行，用户可以重写这个函数实现类似用户验证之类。
func (this *BaseController) Prepare() {
	this.pageSize = 20
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)

	this.auth(controllerName, actionName)

	this.Data["version"] = beego.AppConfig.String("version")
	this.Data["siteName"] = beego.AppConfig.String("site.name")
	this.Data["curRoute"] = this.controllerName + "." + this.actionName
	this.Data["curController"] = this.controllerName
	this.Data["curAction"] = this.actionName
	this.Data["loginUserId"] = this.userId
	this.Data["loginUserName"] = this.userName
}

//登录状态验证
func (this *BaseController) auth(controllerName, actionName string) {
	arr := strings.Split(this.Ctx.GetCookie("auth"), "|")
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			user, err := models.UserGetById(userId)
			if err == nil && password == libs.Md5([]byte(this.getClientIp()+"|"+user.Password+user.Salt)) {
				this.userId = user.Id
				this.userName = user.UserName
				this.user = user
			}

			//用户已登录状态  添加用户访问日志
			t := time.Now()
			params := this.Ctx.Request
			fmt.Println(params)

			log := new(BaseModels.Log)
			log.Uid = this.userId
			log.Action = controllerName + "." + actionName //控制器+方法
			log.Ip = this.getClientIp() //访问者的ip
			//log.Params = this.Ctx.Request.Header; //请求参数
			log.Type = 0
			log.CreateTime = t.Unix()
			BaseModels.LogAdd(log)
		}
	}
	//未登录重定向至登录界面
	if this.userId == 0 && (this.controllerName != "login" ||
		(this.controllerName == "login" && this.actionName != "logout" && this.actionName != "login")) {
		this.redirect(beego.URLFor("AdminController.Login"))
	}
}

//渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + this.controllerName + "/" + this.actionName + ".html"
	}
	this.Layout = "layout/layout.html"
	this.TplName = tplname
}

// 重定向
func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
	this.StopRun()
}

// 是否POST提交
func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}

// 显示错误信息
func (this *BaseController) showMsg(args ...string) {
	this.Data["message"] = args[0]
	redirect := this.Ctx.Request.Referer()
	if len(args) > 1 {
		redirect = args[1]
	}

	this.Data["redirect"] = redirect
	this.Data["pageTitle"] = "系统提示"
	this.display("error/message")
	this.Render()
	this.StopRun()
}

// 输出json
func (this *BaseController) jsonResult(out interface{}) {
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["msg"] = msg

	this.jsonResult(out)
}

//获取用户IP地址
func (this *BaseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}
