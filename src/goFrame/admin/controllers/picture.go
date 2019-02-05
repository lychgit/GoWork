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

	//fileKey := "filename"
	//pathSeparator := (string)(os.PathSeparator) //目录分隔符
	cleanupTargetDir := true // Remove old files
	//maxFileAge := 5 * 3600 // Temp file age in seconds
	//beego.Debug(maxFileAge)

	tempPath := "upload/tmp"
	uploadDir := "upload"
	// Create target dir  判断保存文件的文件夹是否存在  不存在则新建
	if k := utils.IsExist(uploadDir); !k {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Create target dir 判断存储上传文件缓存的文件夹是否存在  不存在则新建
	if k := utils.IsExist(tempPath); !k {
		os.Mkdir(tempPath, os.ModePerm)
	}

	// Remove old temp files
	if (cleanupTargetDir) {
		if v := utils.IsDir(tempPath); !v {
			//'{"jsonrpc" : "2.0", "error" : {"code": 100, "message": "Failed to open temp directory."}, "id" : "id"}'
		}
		//递归移除 缓存目录下的旧文件
	}

	beego.Debug("param")
	param := this.GetString("param")
	beego.Debug(param)
	if param == "settask" {
		//保存上传任务信息     保存的文件信息是整个文件的数据信息
		this.saveTaskInfo(tempPath)
	} else if param == "checkchunk" {
		//检测分片是否存在
		this.checkChunk(tempPath)
	} else if param == "mergechunks" {
		//合并文件分片
		this.mergeBlock(nil)
	} else {
		//上传图片
		this.uploadfile(tempPath)
	}

	//file, fileHead, fileErr := this.Ctx.Request.FormFile(fileKey)  //上传的文件
	//beego.Debug(file)
	//beego.Debug(fileHead)
	//if file == nil {
	//	beego.Debug("xx上传失败")
	//	msg := ""
	//	this.jsonResult(enums.JRCodeFailed, msg, nil)
	//}
	//
	//// Get a file name
	//var fileName string
	//if name := this.Ctx.Request.Form.Get("name"); name != "" {
	//	fileName = name
	//} else if  fileErr == nil {
	//	fileName = fileHead.Filename
	//} else {
	//	fileName = utils.UniqueId()//生成一个唯一ID
	//}
	//filePath := tempPath + pathSeparator + fileName //文件缓存路径
	//beego.Debug(filePath)
	//uploadPath := uploadDir + pathSeparator + fileName //文件存储路径
	//beego.Debug(uploadPath)

	//if ($this->request->hasFiles()) {
	//	$upload = new Uploader([
	//		'rootPath' => $this->config['module']['gallery']['uploadDir'],
	//		'savePath' => '/tmp/',
	//	'autoSub' => false,
	//		'saveName' => $taskid . '.' . $chunk,
	//		'saveExt' => 'tmp'
    //        ]);
	//	foreach ($this->request->getUploadedFiles() as $file) {
	//		try {
	//			$data = $upload->upload($file);
	//		} catch (\Exception $e) {
	//			$this->di->get('logger')->error($e->getMessage());
	//			return new \Xin\Lib\MessageResponse('Upload Fail', 'error', [], 500);
	//		}
	//	}
	//	return;
	//}
	//
	//return new \Xin\Lib\MessageResponse("File is Pending", 'error', [], 500);
}

/**
 * 	保存文件上传的任务信息
 *	params string tempPath  存储缓存路径
 */
func (this *PictureController) saveTaskInfo(tempPath string) {
	beego.Debug("saveTaskInfo: " + tempPath)
	pathSeparator := string(os.PathSeparator)
	fileHash := strings.TrimSpace(this.GetString("filehash"))
	chunkSize := strings.TrimSpace(this.GetString("chunksize"))
	taskId := utils.Md5(fileHash + "_" + string(this.userId))
	infoPath := tempPath + pathSeparator + taskId + "info"
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
 *	params string tempPath  分片存储缓存路径
 */
func (this *PictureController) checkChunk(tempPath string) {
	chunk := this.GetString("chunk", "0")
	if chunkSize, ok := interface{}(this.GetString("chunksize", "0")).(int64); !ok {
		taskId := this.GetString("taskid")
		if !this.Empty(chunk) || !this.Empty(chunkSize){
			this.jsonResult(enums.JRCodeFailed, "invalid param", nil)
		}
		data := make(map[string]bool)
		//$isExist = filesize($tmpfile) == $chunkSize;
		tempFile := tempPath + taskId + chunk + ".tmp"
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
func (this *PictureController) mergeBlock(fileFlag interface{}){
	this.jsonResult(enums.JRCodeSucc, "", nil)
}

/**
 * 	上传文件    上传文件生成缓存
 */
func (this *PictureController) uploadfile(filePath string)  {
	beego.Debug("uploadfile 上传文件生成缓存")
	if hasFiles := this.Ctx.Request.ParseMultipartForm(32 << 20); hasFiles != nil {
		this.jsonResult(enums.JRCodeFailed, "上传文件解析失败", nil)
	}
	var fileHeads []*multipart.FileHeader
	fileHeads = this.Ctx.Request.MultipartForm.File["file"]  //获取上传的文件句柄   type: array
	var file multipart.File  //打开的文件句柄
	var err error
	defer file.Close()
	for _, fileHead := range fileHeads {
		beego.Debug(fileHead.Header)
		beego.Debug(fileHead.Filename)  //文件名称
		beego.Debug(fileHead.Size)	//文件大小
		file, err = fileHead.Open()
		if err != nil {
			beego.Debug("文件打开失败")
			this.jsonResult(enums.JRCodeFailed, err.Error(), nil)
		} else {
			beego.Debug("文件打开成功")

		}
	}
	this.jsonResult(enums.JRCodeSucc, "", nil)
}