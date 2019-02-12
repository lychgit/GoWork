package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/enums"
	"goFrame/utils"
	)

type TestController struct {
	BaseController
}

func (this *TestController) Index() {

	beego.Debug("init - upload")
	uploadConf := make(map[string]interface{})
	uploadConf["MaxSize"] = 100
	uploadConf["AutoSub"] = true
	utils.NewUpload(uploadConf)
	this.jsonResult(enums.JRCodeFailed,"",nil)
}
