package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

// BackendUserQueryParam 用于查询的类
type UserQueryParam struct {
	BaseQueryParam
	UserNameLike string //模糊查询
	MobileLike   string //精确查询
	SearchStatus string //为空不查询，有值精确查询
}

type User struct {
	Id             int `orm:"auto"`
	UserName       string
	Password       string
	Salt           string
	LastLogin      int64
	LastIp         string
	Status         int
	IsSuper        bool
	Mobile         string         `orm:"size(16)"`
	Email          string         `orm:"size(256)"`
	Avatar         string         `orm:"size(256)"`
	RoleIds        []int          `orm:"-" form:"RoleIds"`
	RoleUserRel    []*RoleUserRel `orm:"reverse(many)"` // 设置一对多的反向关系
	MenuUrlForList []string       `orm:"-"`
	CreateTime     int64
}

func (u *User) TableName() string {
	return TableName("user")
}

func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *User) Add() error {
	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}

func UserAdd(user *User) (int64, error) {
	return orm.NewOrm().Insert(user)
}

func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}

func UserGetById(id int) (*User, error) {
	u := new(User)
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserGetByName(userName string) (*User, error) {
	u := new(User)
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("user_name", userName).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserList(page, pageSize int, filters ...interface{}) ([]*User, int64) {
	offset := (page - 1) * pageSize
	users := make([]*User, 0)
	query := orm.NewOrm().QueryTable(TableName("user"))
	if len(filters) > 0 {
		l := len(filters)
		for i := 0; i < l; i += 2 {
			query = query.Filter(filters[i].(string), filters[i+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&users)
	return users, total
}


func UserPageList(params *UserQueryParam) ([]*User, int64) {
	query := orm.NewOrm().QueryTable(TableName("user"))
	data := make([]*User, 0)
	//默认排序
	sortOrder := "Id"
	switch params.Sort {
	case "Id":
		sortOrder = "Id"
	}
	if params.Order == "desc" {
		sortOrder = "-" + sortOrder
	}
	query = query.Filter("username__istartswith", params.UserNameLike)
	query = query.Filter("mobile__istartswith", params.MobileLike)
	//if len(params.Mobile) > 0 {
	//	query = query.Filter("mobile", params.Mobile)
	//}
	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}
	total, _ := query.Count()
	query.OrderBy(sortOrder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}


func UserListGrid(page, pageSize int, filters ...interface{}) []User {
	data, total := UserList(page, pageSize)
	list := make([]User, total)
	for i, item := range data {
		list[i] = *item
	}
	return list
}
