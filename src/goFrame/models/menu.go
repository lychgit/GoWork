package models

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	)

type Menu struct {
	Id           int `orm: "auto"`
	Pid          int    //上一级菜单id
	TitleCn      string //菜单中文名称
	TitleEn      string //菜单英文名称
	Icon         string //菜单图标
	Sort         int    //排序
	Type         int    //类型
	LogicDelete  int    //逻辑删除
	CreateTime   int64  //创建时间
	UpdateTime   int64  //更新时间
	LinkUrl      string         `orm:"-"`
	UrlFor       string         `orm:"size(256)" Json:"-"` //菜单url链接
	HtmlDisabled int            `orm:"-"`                  //在html里应用时是否可用
	Level        int            `orm:"-"`                  //第几级，从0开始
	Parent       *Menu          `orm:"null;rel(fk)"`       // RelForeignKey relation
	Sons         []*Menu        `orm:"reverse(many)"`      // fk 的反向关系
	SonNum       int            `orm:"-"`
	RoleMenuRel  []*RoleMenuRel `orm:"reverse(many)"` // 设置一对多的反向关系
}

func (c *Menu) TableName() string {
	return TableName("menu")
}

func (c *Menu) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func (c *Menu) Add() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}
	return nil
}

func (c *Menu) Delete() error {
	if _, err := orm.NewOrm().Delete(c); err != nil {
		return err
	}
	return nil
}

func MenuGetById(id int) (*Menu, error) {
	c := new(Menu)
	err := orm.NewOrm().QueryTable(TableName("menu")).Filter("id", id).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func MenuGetByName(name string) (*Menu, error) {
	c := new(Menu)
	err := orm.NewOrm().QueryTable(TableName("menu")).Filter("name", name).One(c)
	if err != nil {
		return nil, err
	}
	return c, err
}

func MenuList(page, pageSize int, filters ...interface{}) ([]*Menu, int64) {
	offset := (page - 1) * pageSize
	menus := make([]*Menu, 0)
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

/*
	根据用户id获取有权管理的菜单列表，并整理成teegrid格式
	先从cache中获取menu列表

	使用 Raw SQL 查询，无需使用 ORM 表定义
	多数据库，都可直接使用占位符号 ?，自动转换
	查询时的参数，支持使用 Model Struct 和 Slice, Array
	o := orm.NewOrm()
	o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
 */
func MenuListGetByUid(uid, maxrtype int) []*Menu {
	var list []*Menu
	//缓存key
	//cachekey := fmt.Sprintf("rms_MenuListGetByUid_%v_%v", uid, maxrtype)
	//从缓存cache中获取用户可查看的菜单列表
	//if err := utils.GetCache(cachekey, &list); err == nil {
	//	return list
	//}
	user, err := UserGetById(uid)
	if err != nil || user == nil {
		return list
	}
	o := orm.NewOrm()
	var sql string
	if user.Id == 0 {
		//如果是管理员，则查出所有的
		sql = fmt.Sprintf(`SELECT id,  pid, title_cn, title_en, type, icon, sort, url_for FROM %s Where logic_delete = 0 Order By sort asc,Id asc`, TableName("menu"))
		o.Raw(sql).QueryRows(&list)
	} else {
		//联查多张表，找出某用户有权管理的
		sql = fmt.Sprintf(`SELECT DISTINCT T0.role_id, T2.id, T2.pid, T2.title_cn, T2.title_en, T2.type, T2.icon, T2.sort, T2.url_for
		FROM %s As T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T1.menu_id = T2.id
		WHERE T0.uid = ? and T2.logic_delete = 0  Order By T2.sort asc,T2.id asc`, TableName("role_user_rel"),TableName("role_menu_rel"), TableName("menu"))
		o.Raw(sql, user.Id).QueryRows(&list)
	}
	//格式化菜单列表
	result := menuListTreeGrid(list)
	//存入缓存cache
	//utils.SetCache(cachekey, result, 30)
	//beego.Debug(result)
	return result
}

// menuListTreeGrid 将资源列表转成treegrid格式
func menuListTreeGrid(list []*Menu) []*Menu {
	result := make([]*Menu, 0)
	for _, item := range list {
		if item.Parent == nil || item.Parent.Id == 0 {
			item.Level = 0
			result = append(result, item)
			result = menuAddSons(item, list, result)
		}
	}
	return result
}

//menuAddSons 添加子菜单
func menuAddSons(cur *Menu, list, result []*Menu) []*Menu {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			cur.SonNum++
			item.Level = cur.Level + 1
			result = append(result, item)
			result = menuAddSons(item, list, result)
		}
	}
	return result
}
