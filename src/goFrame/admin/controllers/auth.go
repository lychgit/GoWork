package controllers

import (
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Index() {
	//需要权限控制
	//this.checkAuthor()
	//将页面左边菜单的某项激活
	this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	//this.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	//this.Data["canEdit"] = this.checkActionAuthor("MenuController", "Edit")
	//this.Data["canDelete"] = this.checkActionAuthor("MenuController", "Delete")
	this.display()
}

