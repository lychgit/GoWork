package models

import "github.com/astaxie/beego/orm"

//日志

type Picture struct {
	Id         int `orm:"auto"`
	Title      string //图片名称
	Path       string //图片路径
	Hash       string //图片哈希值
	Size       string //图片大小
	Status     int64  //图片状态
	Type       string //图片格式
	CreateTime int64  //图片生成时间
}

func (p *Picture) TableName() string {
	return TableName("picture")
}

//add
func PictureAdd(p *Picture) (int64, error) {
	return orm.NewOrm().Insert(p)
}

//delete
func PictureDeleteById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("picture")).Filter("id", id).Delete()
	return err
}

//find
func PictureFindById(id int) (*Picture, error) {
	picture := &Picture{
		Id: id,
	}
	err := orm.NewOrm().Read(picture)
	if err != nil {
		return nil, err

	}
	return picture, nil
}

func PictureGetList(page, pageSize int, filters ...interface{}) ([] *Picture, int64) {
	offset := pageSize * (page - 1)
	pictures := make([] *Picture, 0)

	query := orm.NewOrm().QueryTable("picture")
	if len(filters) > 0 {
		length := len(filters)
		for k := 0; k < length; k += 2 {
			//add condition expression to QuerySeter. for example:  filter by UserName == 'slene'
			query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count();
	// OrderBy:    "column" means ASC, "-column" means DESC.
	// Limit: 		qs.Limit(10, 2)    sql-> limit 10 offset 2
	query.OrderBy("-id").Limit(pageSize, offset).All(&pictures)
	return pictures, total
}
