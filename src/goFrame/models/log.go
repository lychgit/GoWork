package models

import "github.com/astaxie/beego/orm"

//日志

type Log struct {
	Id          	int			`orm:"auto"`
	Uid				int			//用户id
	Action      	string		//访问的方法
	Ip				string		//访问ip
	Params			string		//请求参数
	Error       	string		//错误信息
	CreateTime  	int64		//日志生成时间
	Type			int64		//日志类型
}

func (l *Log) TableName() string {
	return TableName("log")
}

//add
func LogAdd(l *Log) (int64, error) {
	return orm.NewOrm().Insert(l)
}

//delete
func LogDeleteById(id int) error {
	_ , err := orm.NewOrm().QueryTable(TableName("log")).Filter("id", id).Delete()
	return err
}

//find
func LogFindById(id int) (* Log, error)  {
	log := &Log{
		Id: id,
	}
	err := orm.NewOrm().Read(log)
	if err != nil {
		return nil, err

	}
	return log, nil
}

func LogGetList(page, pageSize int, filters ...interface{}) ([] *Log, int64) {
	offset := pageSize * (page - 1)
	logs := make([] *Log, 0)

	query := orm.NewOrm().QueryTable("log")
	if len(filters) > 0 {
		length := len(filters)
		for k := 0; k < length; k+=2 {
			//add condition expression to QuerySeter. for example:  filter by UserName == 'slene'
			query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count();
	// OrderBy:    "column" means ASC, "-column" means DESC.
	// Limit: 		qs.Limit(10, 2)    sql-> limit 10 offset 2
	query.OrderBy("-id").Limit(pageSize, offset).All(&logs)
	return logs, total
}

/**
展示时写法
// 最近执行失败的日志
	logs, _ = models.TaskLogGetList(1, 20, "status__lt", 0)
	errLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, err := models.TaskGetById(v.TaskId)
		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["error"] = beego.Substr(v.Error, 0, 100)
		row["status"] = v.Status
		errLogs[k] = row
	}

	this.Data["recentLogs"] = recentLogs
	this.Data["errLogs"] = errLogs
 */