package controllers

import (
	"goFrame/models"
	"goFrame/enums"
	"github.com/astaxie/beego/orm"
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

//DataList 角色列表
func (this *RoleController) RoleList() {
	//获取角色列表
	data := models.RoleListGrid(this.page, this.pageSize)
	//定义返回的数据结构
	this.jsonResult(enums.JRCodeSucc, "", data)
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

func (this *RoleController) Save() {
	m := models.Role{}
	o := orm.NewOrm()
	var err error
	//获取form里的值
	if err = this.ParseForm(&m); err != nil {
		this.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	if m.Id == 0 {
		if _, err := o.Insert(&m); err != nil {
			this.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		//if role, err := models.RoleGetById(m.Id); err != nil {
		//	this.jsonResult(enums.JRCode004, "数据无效，请刷新后重试", m.Id)
		//} else {
		//
		//}
		if _, err := o.Update(&m); err != nil {
			this.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}

func (this *RoleController) RoleEdit() {
	// Edit 添加 编辑 页面
	//如果是Post请求，则由Save处理
	if this.Ctx.Request.Method == "POST" {
		this.Save()
	}
	Id, _ := this.GetInt(":id", 0)
	role := &models.Role{}
	var err error
	if Id > 0 {
		role, err = models.RoleGetById(Id)
		if err != nil {
			this.pageError("数据无效，请刷新后重试")
		}
		//o := orm.NewOrm()
		//o.LoadRelated(role, "RoleUserRel")
	} else {
		//添加用户时默认状态为启用
		//role.Status = enums.Enabled
	}
	this.Data["m"] = role

	this.setTpl("admin/role/edit.html", "layout/admin/layout_pullbox.html")
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["footerjs"] = "admin/user/edit_js.html"
}
