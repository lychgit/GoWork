package controllers

import (
	"goFrame/models/system"
	"github.com/astaxie/beego"
	"time"
)

type ConfigController struct {
	BaseController
}

func (this *ConfigController) Config() {
	this.Data["pageTitle"] = "系统配置"
	configs, _ := system.ConfigList(1, 20, "status", 1)
	configList := make([]map[string]interface{}, len(configs))
	for k, v := range configs {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["desc"] = v.Desc
		row["settings"] = v.Settings
		row["status"] = v.Status
		row["type"] = v.Type
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		configList[k] = row
	}
	this.Data["configList"] = configList
	this.display()
}

func (this *ConfigController) EditConfig() {
	if this.isPost() {
		out := AjaxJson{}
		out.status = true
		out.data = make(map[string]string)
		out.data["errorMsg"] = ""
		name := this.Ctx.Input.Query("name")
		settings := this.Ctx.Input.Query("settings")
		desc := this.Ctx.Input.Query("desc")
		status := this.Ctx.Input.Query("status")
		sort := this.Ctx.Input.Query("sort")
		if name == "" {
			out.data["errorMsg"] = "配置名称不能为空"
		}
		if settings == "" {
			out.data["errorMsg"] = "配置值不能为空"
		}
		if desc == "" {
			out.data["errorMsg"] = "配置说明不能为空"
		}
		if out.data["errorMsg"] == "" {
			t := time.Now()
			c := new(system.Config)
			c.Name = name
			c.Settings = settings
			c.Desc = desc
			if value, ok := interface{}(status).(int); ok {
				c.Status = value
			} else {
				c.Status = 0
			}
			if value, ok := interface{}(sort).(int); ok {
				c.Sort = value
			} else {
				c.Sort = 0
			}
			c.Type = 0
			c.CreateTime = t.Unix()
			c.UpdateTime = t.Unix()
			if err := c.Add(); err != nil {
				out.status = false
				out.data["errorMsg"] = err.Error()
			}
		} else {
			out.status = false
		}
		this.jsonResult(interface{}(!out.status).(int), out.data["errorMsg"], out.data)
	}
	this.display()
}

func (this *ConfigController) DeleteConfig() {
	if this.isPost() {
		out := AjaxJson{}
		out.status = true
		out.data = make(map[string]string)
		out.data["errorMsg"] = ""
		//id := this.Ctx.Input.Query("id")
		id, _ := this.GetInt("id")
		if config, err := system.ConfigGetById(id); err != nil {
			out.status = false
			out.data["errorMsg"] = "can't found this config"
		} else {
			if err := config.Delete(); err != nil {
				out.status = false
				out.data["errorMsg"] = err.Error()
			}
		}
		this.jsonResult(interface{}(!out.status).(int), out.data["errorMsg"], out.data)
	}
}
