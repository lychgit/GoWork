package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/utils"
	"os"
)

type PictureController struct {
	BaseController
}

func (this *PictureController) Index() {
	beego.Debug("PictureController-Index")
	//是否显示更多查询条件的按钮
	this.Data["showMoreQuery"] = true
	//需要权限控制
	this.checkAuthor()
	//将页面左边菜单的某项激活
	//this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	//this.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	//this.Data["canEdit"] = this.checkActionPictureor("MenuController", "Edit")
	//this.Data["canDelete"] = this.checkActionPictureor("MenuController", "Delete")
	this.display()
}
func (this *PictureController) PictureUpload() {
	beego.Debug("PictureController-PictureUpload")
	// Support CORS
	// header("Access-Control-Allow-Origin: *");
	// other CORS headers if any...
	if (this.Ctx.Request.Method == "OPTIONS") {
		return; // finish preflight CORS requests here
	}

	//设置脚本执行时间
	// 5 minutes execution time
	//@set_time_limit(5 * 60);

	targetDir := "upload_tmp"
	uploadDir := "upload"
	fileKey := "request请求中文件名"

	cleanupTargetDir := true // Remove old files
	maxFileAge := 5 * 3600 // Temp file age in seconds

	beego.Debug(cleanupTargetDir)
	beego.Debug(maxFileAge)

	// Create target dir 判断存储上传图片缓存的文件夹是否存在  不存在则新建
	if (!utils.IsExist(targetDir)) {
		os.Mkdir(targetDir, os.ModePerm)
	}
	// Create target dir  判断保存图片的文件夹是否存在  不存在则新建
	if (!utils.IsExist(uploadDir)) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	file, fileHead, fileErr := this.Ctx.Request.FormFile(fileKey)  //上传的文件
	beego.Debug(file)
	beego.Debug(fileHead)

	// Get a file name
	var fileName string
	if name := this.Ctx.Request.Form.Get("name"); name != "" {
		fileName = name
	} else if  fileErr == nil {
		fileName = fileHead.Filename
	} else {
		fileName = utils.UniqueId()//生成一个唯一ID
	}
	spl := (string)(os.PathSeparator)
	filePath := targetDir + spl + fileName //文件缓存路径
	uploadPath := uploadDir + spl + fileName //文件存储路径
	beego.Debug(filePath)
	beego.Debug(uploadPath)

	this.Data["json"] = nil
	this.ServeJSON()
}
