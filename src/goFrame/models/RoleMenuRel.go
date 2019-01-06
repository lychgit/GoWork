package models

import (
	"github.com/astaxie/beego/orm"
)

type RoleMenuRel struct {
	Id         int `orm: "auto"`
	//RoleId     int                    //外键//角色id
	//MenuId     int                    //菜单id
	Role       *Role `orm:"rel(fk)"`  //外键
	Menu       *Menu `orm:"rel(fk)" ` // 外键
	CreateTime int64                  //创建时间
}

func (c *RoleMenuRel) TableName() string {
	return TableName("role_menu_rel")
}

func (c *RoleMenuRel) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *RoleMenuRel) Add() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}
	return nil
}

func (c *RoleMenuRel) Delete() error {
	if _, err := orm.NewOrm().Delete(c); err != nil {
		return err
	}
	return nil
}

func RoleMenuRelGetById(id int) (*RoleMenuRel, error) {
	c := new(RoleMenuRel)
	err := orm.NewOrm().QueryTable(TableName("menu")).Filter("id", id).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func RoleMenuRelGetByName(name string) (*RoleMenuRel, error) {
	c := new(RoleMenuRel)
	err := orm.NewOrm().QueryTable(TableName("menu")).Filter("name", name).One(c)
	if err != nil {
		return nil, err
	}
	return c, err
}

func RoleMenuRelList(page, pageSize int, filters ...interface{}) ([]*RoleMenuRel, int64) {
	offset := (page - 1) * pageSize
	menus := make([]*RoleMenuRel, 0)
	query := orm.NewOrm().QueryTable(TableName("menu"))
	if len(filters) > 0 {
		l := len(filters)
		for i := 0; i < l; i += 2 {
			query = query.Filter(filters[i].(string), filters[i+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-sort", "-id").Limit(pageSize, offset).All(&menus)
	return menus, total
}
