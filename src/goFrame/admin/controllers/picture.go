package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/utils"
	"os"
	"strings"
	"bufio"
	"goFrame/enums"
	"mime/multipart"
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
	if v := this.isPost(); !v {
		this.jsonResult(enums.JRCodeFailed, "request error", nil)
	}
	// Support CORS
	// header("Access-Control-Allow-Origin: *");
	// other CORS headers if any...
	if (this.Ctx.Request.Method == "OPTIONS") {
		return; // finish preflight CORS requests here
	}
	//设置脚本执行时间
	// 5 minutes execution time
	//@set_time_limit(5 * 60);

	pathSeparator := (string)(os.PathSeparator) //目录分隔符
	//pathSeparator := "/" //目录分隔符
	cleanupTargetDir := true // Remove old files
	//maxFileAge := 5 * 3600 // Temp file age in seconds
	//beego.Debug(maxFileAge)

	tmpPath := "tmp"
	rootPath := "upload"
	fileTempPath := tmpPath + pathSeparator//文件缓存路径
	uploadPath := rootPath + pathSeparator//文件存储路径

	// Create target dir  判断保存文件的文件夹是否存在  不存在则新建
	if k := utils.IsExist(rootPath); !k {
		os.Mkdir(rootPath, os.ModePerm)
	}

	// Create target dir 判断存储上传文件缓存的文件夹是否存在  不存在则新建
	if k := utils.IsExist(tmpPath); !k {
		os.Mkdir(tmpPath, os.ModePerm)
	}

	// Remove old temp files
	if (cleanupTargetDir) {
		if v := utils.IsDir(tmpPath); !v {
			beego.Error("Failed to open temp directory!")
			this.jsonResult(enums.JRCodeFailed, "Failed to open temp directory!", nil)
		}
		//递归移除 缓存目录下的旧文件
	}

	param := this.GetString("param")
	beego.Debug("param")
	beego.Debug(param)
	if param == "settask" {
		//保存上传任务信息     保存的文件信息是整个文件的数据信息
		this.saveTaskInfo(tmpPath)
	} else if param == "checkchunk" {
		//检测分片是否存在
		this.checkChunk(tmpPath)
	} else if param == "mergechunks" {
		//合并文件分片
		this.mergeBlock(nil, fileTempPath, uploadPath)
	} else {
		//上传图片
		file, fileHead, fileErr := this.Ctx.Request.FormFile("file")  //上传的文件
		//beego.Debug(file)
		//beego.Debug(fileHead)
		if file == nil {
			beego.Error(" 未找到要上传的图片!, ERROR: FormFile获取上传图片失败!")
			msg := ""
			this.jsonResult(enums.JRCodeFailed, msg, nil)
		}
		//fileBytes, _ := ioutil.ReadAll(file)
		//beego.Debug("fileBytes")
		//beego.Debug(fileBytes)

		// Get a file name
		var fileName string
		if name := this.Ctx.Request.Form.Get("name"); name != "" {
			fileName = name
		} else if  fileErr == nil {
			fileName = fileHead.Filename
		} else {
			fileName = utils.UniqueId()//生成一个唯一ID
		}

		////文件id   上传多个图片 id不同  WU_FILE_1、 WU_FILE_2
		//Id := this.Ctx.Request.Form.Get("id")
		//beego.Debug(Id)
		////文件类型
		//fileType := this.Ctx.Request.Form.Get("type")
		//beego.Debug(fileType) //image/png

		chunk := this.GetString("chunk", "0") //分块下标
		taskId := this.GetString("taskid", "0") //文件唯一标记
		saveName := fileName + "_" +  taskId + "_" + chunk
		this.uploadfile(uploadPath, fileTempPath, saveName)
	}
}

/**
 * 	保存文件上传的任务信息
 *	params string tmpPath  存储缓存路径
 */
func (this *PictureController) saveTaskInfo(tmpPath string) {
	beego.Debug("saveTaskInfo: " + tmpPath)
	pathSeparator := string(os.PathSeparator)
	fileHash := strings.TrimSpace(this.GetString("filehash"))
	chunkSize := strings.TrimSpace(this.GetString("chunksize"))
	taskId := utils.Md5(fileHash + "_" + string(this.userId))
	infoPath := tmpPath + pathSeparator + taskId + "info"
	//判断分片任务信息缓存文件是否存在
	if v := utils.IsFile(infoPath); !v {
		data := make(map[string]interface{})
		data["filename"] = this.GetString("filename")
		data["filehash"] = fileHash
		data["filesize"] = this.GetString("filesize")
		data["label"] = this.GetString("label")
		data["chunksize"] = chunkSize
		data["fileid"] = this.GetString("fileid")
		//file_put_contents($infoPath, serialize($data)); //将任务信息写入infoPath目录下保存
		infoFile, err := os.OpenFile(infoPath, os.O_CREATE|os.O_WRONLY, 0644)
		defer infoFile.Close()
		if err != nil {
			//分片任务信息缓存文件创建失败
			beego.Error("create file error:", err)
		}
		ioW := bufio.NewWriter(infoFile) //创建新的 Writer 对象
		_, error := ioW.WriteString(this.Json_encode(data))
		if error != nil {
			beego.Error("write error", error)
		}
		ioW.Flush()
	}
	var data = make(map[string]interface{});
	data["taskid"] = taskId;
	beego.Debug(data)
	this.jsonResult(enums.JRCodeSucc, "", data)
}

/**
 * 	检测分片是否存在
 *	params string tmpPath  分片存储缓存路径
 *	chunk 分块下标
 *	chunksize 分块大小
 */
func (this *PictureController) checkChunk(tmpPath string) {
	chunk := this.GetString("chunk", "0") //分块下标
	if chunkSize, ok := interface{}(this.GetString("chunksize", "0")).(int64); !ok {
		taskId := this.GetString("taskid")
		if !this.Empty(chunk) || !this.Empty(chunkSize){
			beego.Error("分块下标或分块大小不能为空! ERROR: chunk or chunkSize is empty!")
			this.jsonResult(enums.JRCodeFailed, "invalid param", nil)
		}
		data := make(map[string]bool)
		//$isExist = filesize($tmpfile) == $chunkSize;
		tempFile := tmpPath + taskId + chunk + ".tmp"
		if !utils.IsFile(tempFile) || utils.GetFile(tempFile).Size() == chunkSize {
			data["isExist"] = false
		} else {
			data["isExist"] = true
		}
		this.jsonResult(enums.JRCodeSucc, "", data)
	} else {
		this.jsonResult(enums.JRCodeFailed, "param error", nil)
	}
}

/**
 * 	合并分片
 *	params interface{} fileFlag  可识别属于同一文件的分片标识
 */
func (this *PictureController) mergeBlock(fileFlag interface{}, tmpPath, uploadPath string){
	taskId := this.GetString("taskid", "0") //文件唯一标记
	infoPath := tmpPath + taskId + "info"
	if utils.IsFile(infoPath) {
		//data := ; //获取文件中数据
		//saveDir := time.Date()


		this.jsonResult(enums.JRCodeSucc, "", nil)
	} else {
		beego.Error("Merge Sharding Failed!")
		this.jsonResult(enums.JRCodeFailed, "Merge Sharding Failed", nil)
	}
}

/**
 * 	上传文件    上传文件生成缓存
 */
func (this *PictureController) uploadfile(rootPath, tmpPath, saveName string)  {
	beego.Debug("uploadfile 上传文件生成缓存")
	if hasFiles := this.Ctx.Request.ParseMultipartForm(32 << 20); hasFiles != nil {
		this.jsonResult(enums.JRCodeFailed, "上传文件解析失败", nil)
	}

	var fileHeads []*multipart.FileHeader
	fileHeads = this.Ctx.Request.MultipartForm.File["file"]  //获取上传的文件句柄   type: array
	//上传upload类初始化
	uploadConf := make(map[string]interface{})
	uploadConf["RootPath"] = rootPath //图库根路径
	uploadConf["SavePath"] = tmpPath //分块缓存文件存储路径
	uploadConf["AutoSub"] = false
	uploadConf["SaveName"] = saveName //缓存文件保存名称
	uploadConf["SaveExt"] = ".tmp" //缓存文件后缀
	if upload, err := utils.NewUpload(uploadConf); err == nil {
		var infos [] map[string]interface{}
		for _, fileHead := range fileHeads {
			beego.Debug("upload temp file")
			beego.Debug(upload)
			if info, err := upload.Upload(fileHead); err != nil {
				beego.Debug(err.Error())
				this.jsonResult(enums.JRCodeFailed, fileHead.Filename + "upload failed!", nil)
			} else {
				beego.Debug(info)
				infos = append(infos, info)
			}
		}
		this.jsonResult(enums.JRCodeSucc, "upload success!", infos)
	} else {
		beego.Error("文件上传类实例化失败！ ERROR: file upload class create error!")
		this.jsonResult(enums.JRCodeFailed, "upload create error!", nil)
	}
}