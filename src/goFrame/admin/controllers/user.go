package controllers

import (
	"goFrame/models"
		"github.com/astaxie/beego"
	"github.com/lhtzbj12/sdrms/enums"
	"strconv"
	"strings"
	"github.com/lhtzbj12/sdrms/utils"
	"github.com/astaxie/beego/orm"
	"encoding/json"
)

type UserController struct {
	BaseController
}

func (this *UserController) Index() {
	beego.Debug("UserController-Index")
	//是否显示更多查询条件的按钮
	this.Data["showMoreQuery"] = true
	//需要权限控制
	this.checkAuthor()
	//将页面左边菜单的某项激活
	//this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)
	//获取角色列表
	data := models.RoleListGrid(this.page, this.pageSize)
	//rolelist["rows"] = this.Json_encode(data)
	this.Data["role_rows"] = data
	//页面里按钮权限控制
	//this.Data["canEdit"] = this.checkActionUseror("MenuController", "Edit")
	//this.Data["canDelete"] = this.checkActionUseror("MenuController", "Delete")

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["datagrid"] = "admin/" + this.controllerName + "/datagrid.html"
	this.display()
}


// DataGrid 后台用户管理页 表格获取数据
func (this *UserController) UserDataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.UserQueryParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	beego.Debug(params)
	//获取数据列表和总数
	data, total := models.UserList(this.page, this.pageSize)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) Save() {
	m := models.User{}
	o := orm.NewOrm()
	var err error
	//获取form里的值
	if err = this.ParseForm(&m); err != nil {
		this.jsonResult(enums.JRCode004, "获取数据失败", m.Id)
	}
	//删除已关联的历史数据
	//if _, err := o.QueryTable(models.RoleBackendUserRelTBName()).Filter("backenduser__id", m.Id).Delete(); err != nil {
	//	this.jsonResult(202, "删除历史关系失败", "")
	//}
	if m.Id == 0 {
		//对密码进行加密
		m.Password = utils.String2md5(m.Password)
		if _, err := o.Insert(&m); err != nil {
			this.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if user, err := models.UserGetById(m.Id); err != nil {
			this.jsonResult(enums.JRCode004, "数据无效，请刷新后重试", m.Id)
		} else {
			m.Password = strings.TrimSpace(m.Password)
			if len(m.Password) == 0 {
				//如果密码为空则不修改
				m.Password = user.Password
			} else {
				m.Password = utils.String2md5(m.Password)
			}
			//本页面不修改头像和密码，直接将值附给新m
			m.Avatar = user.Avatar
		}
		if _, err := o.Update(&m); err != nil {
			this.jsonResult(enums.JRCode003, "编辑失败", m.Id)
		}
	}
	////添加关系
	//var relations []models.RoleBackendUserRel
	//for _, roleId := range m.RoleIds {
	//	r := models.Role{Id: roleId}
	//	relation := models.RoleBackendUserRel{BackendUser: &m, Role: &r}
	//	relations = append(relations, relation)
	//}
	//if len(relations) > 0 {
	//	//批量添加
	//	if _, err := o.InsertMulti(len(relations), relations); err == nil {
	//		this.jsonResult(101, "保存成功", m.Id)
	//	} else {
	//		this.jsonResult(202, "保存失败", m.Id)
	//	}
	//} else {
	//	this.jsonResult(101, "保存成功", m.Id)
	//}
}

func (this *UserController) UserEdit() {
	beego.Debug("UserEdit")
	// Edit 添加 编辑 页面
	//如果是Post请求，则由Save处理
	if this.Ctx.Request.Method == "POST" {
		this.Save()
	}
	Id, _ := this.GetInt(":id", 0)
	user := &models.User{}
	var err error
	if Id > 0 {
		user, err = models.UserGetById(Id)
		if err != nil {
			this.pageError("数据无效，请刷新后重试")
		}
		o := orm.NewOrm()
		o.LoadRelated(user, "RoleUserRel")
	} else {
		//添加用户时默认状态为启用
		user.Status = enums.Enabled
	}
	this.Data["m"] = user
	//获取关联的roleId列表
	var roleIds []string
	for _, item := range user.RoleUserRel {
		roleIds = append(roleIds, strconv.Itoa(item.Role.Id))
	}
	this.Data["roles"] = strings.Join(roleIds, ",")

	beego.Debug(this.Data["roles"])

	this.setTpl("admin/user/edit.html", "layout/admin/layout_pullbox.html")
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["footerjs"] = "admin/user/edit_js.html"
}