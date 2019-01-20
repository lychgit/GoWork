package controllers

import (
	"goFrame/models"
	"strings"
	)

type MenuController struct {
	BaseController
}

func (this *MenuController) Index() {
	//需要权限控制
	//this.checkAuthor()
	//将页面左边菜单的某项激活
	this.display()
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	//this.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	//this.Data["canEdit"] = this.checkActionAuthor("MenuController", "Edit")
	//this.Data["canDelete"] = this.checkActionAuthor("MenuController", "Delete")
}

//UserMenuTree 获取用户有权管理的菜单、区域列表
func (this *MenuController) UserMenuTree() {
	uid := this.curUser.Id
	//获取用户有权管理的菜单列表（包括区域）
	tree := models.MenuListGetByUid(uid, 1)
	////转换UrlFor 2 LinkUrl
	this.UrlForLink(tree)
	//beego.Debug(tree)
	this.jsonResult(0, "", tree)
}

//UrlForLink 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (this *MenuController) UrlForLink(menus []*models.Menu) {
	for _, item := range menus {
		//beego.Debug(item.UrlFor)
		item.LinkUrl = this.UrlForLinkOne(item.UrlFor)
		//beego.Debug(item.LinkUrl)
	}
}

// UrlForLinkOne 使用URLFor方法，将资源表里的UrlFor值转成LinkUrl
func (this *MenuController) UrlForLinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}
	// MenuController.Edit,:id,1
	strs := strings.Split(urlfor, ",")
	if len(strs) == 1 {
		return this.URLFor(strs[0])
	} else if len(strs) > 1 {
		var values []interface{}
		for _, val := range strs[1:] {
			values = append(values, val)
		}
		return this.URLFor(strs[0], values...)
	}
	return ""
}


//
////TreeGrid 获取所有资源的列表
//func (this *MenuController) TreeGrid() {
//	tree := models.MenuTreeGrid()
//	//转换UrlFor 2 LinkUrl
//	this.UrlForLink(tree)
//	//this.jsonResult(0, "", tree)
//}


//ParentTreeGrid 获取可以成为某节点的父节点列表
//func (this *MenuController) ParentTreeGrid() {
//	Id, _ := this.GetInt("id", 0)
//	tree := models.MenuTreeGrid4Parent(Id)
//	//转换UrlFor 2 LinkUrl
//	this.UrlForLink(tree)
//	//this.jsonResult(0, "", tree)
//}
