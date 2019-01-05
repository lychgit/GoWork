package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/models"
	"strings"
	"strconv"
	"goFrame/libs"
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	curUser        models.User
	userId         int
	userName       string
	pageSize       int
}

type AjaxJson struct {
	status bool
	data   map[string]string
}

//这个函数主要是为了用户扩展用的，这个函数会在Get、Post、Delete、Put、Finish等这些 Method 方法之前执行，用户可以重写这个函数实现类似用户验证之类。
func (this *BaseController) Prepare() {
	this.pageSize = 20
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)

	//判断用户是否有权访问某地址，无权则会跳转到错误页面
	//this.checkAuthor("login")

	this.getUserSession();
	//判断用户登录状态
	this.auth(controllerName, actionName)

	this.Data["pageTitle"] = "LYCH System Backstage"
	this.Data["version"] = beego.AppConfig.String("version")
	this.Data["siteName"] = beego.AppConfig.String("site.name")
	this.Data["curRoute"] = this.controllerName + "." + this.actionName
	this.Data["curController"] = this.controllerName
	this.Data["curAction"] = this.actionName
	this.Data["loginUserId"] = this.userId
	this.Data["loginUserName"] = this.userName
}

// 传入的参数为忽略权限控制的Action
func (this *BaseController) checkAuthor(ignores ...string) {
	//如果Action在忽略列表里，则直接通用
	for _, ignore := range ignores {
		if ignore == this.actionName {
			return
		}
	}
	hasAuthor := false //c.checkActionAuthor(c.controllerName, c.actionName)
	if !hasAuthor {
		//utils.LogDebug(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", c.controllerName, c.actionName, c.curUser.Id))
		//如果没有权限
		if !hasAuthor {
			if this.Ctx.Input.IsAjax() {
				//c.jsonResult(enums.JRCode401, "无权访问", "")
			} else {
				this.pageError("无权访问")
			}
		}
	}
}

//登录状态验证
func (this *BaseController) auth(controllerName, actionName string) {
	beego.Debug("here is auth")                    //debug埋点
	beego.Debug(controllerName + "." + actionName) //debug埋点
	//beego.Debug(this.Ctx.Request.URL)              //debug埋点
	//beego.Debug(this.Ctx.Request)                  //debug埋点
	arr := strings.Split(this.Ctx.GetCookie("auth"), "|")
	//beego.Debug("GetCookie" + this.Ctx.GetCookie("auth")) //debug埋点
	//beego.Debug(arr[0]) //debug埋点
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			user, err := models.UserGetById(userId)
			if err == nil && password == libs.Md5([]byte(this.getClientIp()+"|"+user.Password+user.Salt)) {
				this.userId = user.Id
				this.userName = user.UserName
				this.curUser = *user
			}
			//用户已登录状态  添加用户访问日志
			//t := time.Now()
			//params := this.Ctx.Request
			//fmt.Println(params)
			//log := new(models.Log)
			//log.Uid = this.userId
			//log.Action = controllerName + "." + actionName //控制器+方法
			//log.Ip = this.getClientIp()                    //访问者的ip
			////log.Params = this.Ctx.Request.Header; //请求参数
			//log.Type = 0
			//log.CreateTime = t.Unix()
			//models.LogAdd(log)
		}
	}
	beego.Debug(this.userId)
	//未登录重定向至登录界面
	if this.userId == 0 && (this.controllerName != "login" ||
		(this.controllerName == "login" && this.actionName != "logout" && this.actionName != "login" && this.actionName != "gettime")) {
		this.redirect(beego.URLFor("LoginController.Login"))
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
	this.Layout = "layout/admin/layout.html"
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["head_resource"] = "template/head.html" //公共css、js等资源文件
	this.TplName = tplname
	//beego.Debug("display:" + this.TplName) //debug埋点
}

// 设置模板
// 第一个参数模板，第二个参数为layout
func (this *BaseController) setTpl(template ...string) {
	var tplName string
	layout := "layout/admin/layout.html"
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//不要Controller这个10个字母
		ctrlName := strings.ToLower(this.controllerName[0 : len(this.controllerName)-10])
		actionName := strings.ToLower(this.actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}
	this.Layout = layout
	this.TplName = tplName
}

// 重定向
func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
	this.StopRun()
}

// 是否POST提交
func (this *BaseController) isPost() bool {
	return strings.ToUpper(this.Ctx.Request.Method) == "POST"
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

// 重定向 去错误页
func (this *BaseController) pageError(msg string) {
	error_url := this.URLFor("AdminController.Error") + "/" + msg
	this.Redirect(error_url, 302)
	this.StopRun()
}

/**
	设置用户session信息
 */
func (this *BaseController) setUserSession(uid int) error {
	//m, err := models.
	m, err := models.UserGetById(uid)
	if err != nil {
		return err
	}
	//获取这个用户能获取到的所有菜单列表
	resourceList := models.MenuListGetByUid(uid, 1000)
	beego.Debug("获取这个用户能获取到的所有菜单列表")
	beego.Debug(resourceList)
	for _, item := range resourceList {
		m.MenuUrlForList = append(m.MenuUrlForList, strings.TrimSpace(item.UrlFor))
	}
	this.SetSession("adminuser", *m)
	return nil
}

func (this *BaseController) getUserSession() {
	a := this.GetSession("adminuser")
	beego.Debug("getUserSession")
	beego.Debug(a)
	if a != nil {
		this.curUser = a.(models.User)
		this.Data["user"] = a
	}
}
