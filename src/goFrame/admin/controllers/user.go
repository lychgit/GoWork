package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"encoding/json"
	"goFrame/enums"
	"goFrame/utils"
	"goFrame/models"
	"time"
	"fmt"
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

/*
	DataGrid 后台用户管理页 表格获取数据
 */
func (this *UserController) UserDataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.UserQueryParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	//data, total := models.UserList(this.page, this.pageSize)
	data, total := models.UserPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	this.Data["json"] = result
	this.ServeJSON()
}

/*
	后台保存用户信息
 */
func (this *UserController) Save() {
	m := models.User{}
	o := orm.NewOrm()
	var err error
	//获取form里的值
	if err = this.ParseForm(&m); err != nil {
		this.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	//删除已关联的历史数据
	//if _, err := o.QueryTable(models.RoleBackendUserRelTBName()).Filter("backenduser__id", m.Id).Delete(); err != nil {
	//	this.jsonResult(202, "删除历史关系失败", "")
	//}
	if m.Id == 0 {
		//新增用户
		m.Password = utils.String2md5(m.Password) //对密码进行加密
		m.CreateTime = time.Now().Unix()
		if _, err := o.Insert(&m); err != nil {
			this.jsonResult(enums.JRCodeFailed, err.Error(), m.Id)
		}
	} else {
		//修改用户信息
		if user, err := models.UserGetById(m.Id); err != nil {
			this.jsonResult(enums.JRCodeFailed, err.Error(), m.Id)
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
			this.jsonResult(enums.JRCodeFailed, err.Error(), m.Id)
		}
	}
	this.jsonResult(enums.JRCodeSucc, "success!", m.Id)
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

/*
	后台编辑用户信息
 */
func (this *UserController) UserEdit() {
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
	this.setTpl("admin/user/edit.html", "layout/admin/layout_pullbox.html")
}

/*
	后台逻辑删除用户
 */
func (this *UserController) UserDelete() {
	if this.Ctx.Request.Method == "POST" {
		userIds := this.GetString("ids")
		ids := make([]int, 0, len(userIds))
		for _, id := range strings.Split(userIds, ",") {
			if id, err := strconv.Atoi(id); err == nil {
				ids = append(ids, id)
			} else {
				this.jsonResult(enums.JRCodeFailed, err.Error(), nil)
			}
		}
		query := orm.NewOrm().QueryTable(models.TableName("user"))
		//物理删除
		//if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		//逻辑删除
		if num, err := query.Filter("id__in", ids).Update(orm.Params{ "logic_delete": 1}); err == nil {
			this.jsonResult(enums.JRCodeSucc, fmt.Sprintf("Successful deletion of %d records", num), nil)
		} else {
			this.jsonResult(enums.JRCodeFailed, "delete failed!", nil)
		}
	} else {
		this.jsonResult(enums.JRCodeSucc, "reuqest method error!", nil)
	}
}