package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"goFrame/models"
	"goFrame/libs"
	)

type AdminController struct {
	BaseController
}

// 首页
func (this *AdminController) Index() {
	//
	//// 即将执行的任务
	//entries := jobs.GetEntries(30)
	//jobList := make([]map[string]interface{}, len(entries))
	//for k, v := range entries {
	//	row := make(map[string]interface{})
	//	job := v.Job.(*jobs.Job)
	//	row["task_id"] = job.GetId()
	//	row["task_name"] = job.GetName()
	//	row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
	//	jobList[k] = row
	//}
	//
	//// 最近执行的日志
	//logs, _ := models.TaskLogGetList(1, 20)
	//recentLogs := make([]map[string]interface{}, len(logs))
	//for k, v := range logs {
	//	task, err := models.TaskGetById(v.TaskId)
	//	taskName := ""
	//	if err == nil {
	//		taskName = task.TaskName
	//	}
	//	row := make(map[string]interface{})
	//	row["task_name"] = taskName
	//	row["id"] = v.Id
	//	row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
	//	row["process_time"] = float64(v.ProcessTime) / 1000
	//	row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
	//	row["output"] = beego.Substr(v.Output, 0, 100)
	//	row["status"] = v.Status
	//	recentLogs[k] = row
	//}
	//
	////// 最近执行失败的日志
	//logs, _ = models.TaskLogGetList(1, 20, "status__lt", 0)
	//errLogs := make([]map[string]interface{}, len(logs))
	//for k, v := range logs {
	//	task, err := models.TaskGetById(v.TaskId)
	//	taskName := ""
	//	if err == nil {
	//		taskName = task.TaskName
	//	}
	//	row := make(map[string]interface{})
	//	row["task_name"] = taskName
	//	row["id"] = v.Id
	//	row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
	//	row["process_time"] = float64(v.ProcessTime) / 1000
	//	row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
	//	row["error"] = beego.Substr(v.Error, 0, 100)
	//	row["status"] = v.Status
	//	errLogs[k] = row
	//}
	//this.Data["recentLogs"] = recentLogs
	//this.Data["errLogs"] = errLogs
	//this.Data["jobs"] = jobList
	//this.Data["cpuNum"] = runtime.NumCPU()
	this.display()
}

// 个人信息
func (this *AdminController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user, _ := models.UserGetById(this.userId)

	if this.isPost() {
		flash := beego.NewFlash()
		user.Email = this.GetString("email")
		user.Update()
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 {
				flash.Error("密码长度必须大于6位")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else if password2 != password1 {
				flash.Error("两次输入的密码不一致")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
				user.Update()
			}
		}
		flash.Success("修改成功！")
		flash.Store(&this.Controller)
		this.redirect(beego.URLFor(".Profile"))
	}

	this.Data["pageTitle"] = "个人信息"
	this.Data["user"] = user
	this.display()
}

