package controllers

import (
	"goFrame/models"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Index() {
	//需要权限控制
	this.checkAuthor()
	//将页面左边菜单的某项激活
	this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	//this.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	//this.Data["canEdit"] = this.checkActionRoleor("MenuController", "Edit")
	//this.Data["canDelete"] = this.checkActionRoleor("MenuController", "Delete")
	this.display()
}

// DataGrid 后台用户管理页 表格获取数据
func (this *RoleController) RoleDataGrid() {
	//直接反序化获取json格式的requestbody里的值
	//var params models.RoleQueryParam
	//json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.RoleList(this.page, this.pageSize)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	this.Data["json"] = result
	this.ServeJSON()
}

