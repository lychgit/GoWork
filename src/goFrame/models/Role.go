package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

// BackendRoleQueryParam 用于查询的类
type RoleQueryParam struct {
	BaseQueryParam
	RoleNameLike string //模糊查询
	SearchStatus string //为空不查询，有值精确查询
}

type Role struct {
	Id          int `orm: "auto"`
	NameCn      string //角色中文名称
	NameEn      string //角色英文名称
	Status      int
	LogicDelete int
	CreateTime  int64 //创建时间
	UpdateTime  int64
}

func (c *Role) TableName() string {
	return TableName("role")
}

func (c *Role) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *Role) Add() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}
	return nil
}

func (c *Role) Delete() error {
	if _, err := orm.NewOrm().Delete(c); err != nil {
		return err
	}
	return nil
}

func RoleGetById(id int) (*Role, error) {
	r := new(Role)
	err := orm.NewOrm().QueryTable(TableName("role")).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func RoleGetByName(name string) (*Role, error) {
	c := new(Role)
	err := orm.NewOrm().QueryTable(TableName("role")).Filter("name", name).One(c)
	if err != nil {
		return nil, err
	}
	return c, err
}

func RoleList(page, pageSize int, filters ...interface{}) ([]*Role, int64) {
	offset := (page - 1) * pageSize
	roles := make([]*Role, 0)
	query := orm.NewOrm().QueryTable(TableName("role"))
	if len(filters) > 0 {
		l := len(filters)
		for i := 0; i < l; i += 2 {
			query = query.Filter(filters[i].(string), filters[i+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&roles)
	return roles, total
}

func RoleListGrid(page, pageSize int, filters ...interface{}) []Role {
	data, total := RoleList(page, pageSize)
	list := make([]Role, total)
	for i, item := range data {
		list[i] = *item
	}
	return list
}

/*
	DataGrid获取Role数据
 */
func RolePageList(params *RoleQueryParam) ([]*Role, int64) {
	beego.Debug("RolePageList")
	query := orm.NewOrm().QueryTable(TableName("role"))
	roles := make([]*Role, 0)
	//默认排序
	sortOrder := "Id"
	switch params.Sort {
	case "Id":
		sortOrder = "Id"
	}
	if params.Order == "desc" {
		sortOrder = "-" + sortOrder
	}
	query = query.Filter("namecn__istartswith", params.RoleNameLike)
	query = query.Filter("logic_delete", 0)
	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}
	total, _ := query.Count()
	query.OrderBy(sortOrder).Limit(params.Limit, params.Offset).All(&roles)
	beego.Debug("roles")
	beego.Debug(roles)
	beego.Debug("total")
	beego.Debug(total)
	return roles, total
}