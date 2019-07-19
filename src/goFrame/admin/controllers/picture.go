package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/utils"
	"os"
	"bufio"
	"goFrame/enums"
	"mime/multipart"
	"time"
	"errors"
	"encoding/json"
	"strings"
)

type PictureController struct {
	BaseController
}

type TempInfo struct {
	ChunkSize int64  `json:"ChunkSize"`
	FileHash  string `json:"FileHash"`
	FileId    string `json:"FileId"`
	FileName  string `json:"FileName"`
	FileSize  int64  `json:"FileSize"`
	Label     string `json:"Label"`
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
	fileTempPath := tmpPath + pathSeparator //文件缓存路径
	uploadPath := rootPath + pathSeparator  //文件存储路径

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

	taskId := this.GetString("taskid", "0") //文件唯一标记

	param := this.GetString("param")
	beego.Debug("param")
	beego.Debug(param)
	if param == "settask" {
		//保存上传任务信息     保存的文件信息是整个文件的数据信息
		fileHash := strings.TrimSpace(this.GetString("filehash"))
		fileId := this.GetString("fileid")
		fileName := this.GetString("filename")
		fileSize, _ := this.GetInt64("filesize", 0)
		chunkSize, _ := this.GetInt64("chunksize", 0)
		label := this.GetString("label")
		data, err := utils.SaveTaskInfo(tmpPath, fileHash, fileId, fileName, label, chunkSize, fileSize, int64(this.userId))
		if err != nil {
			this.jsonResult(enums.JRCodeFailed, "", data)
		}
		this.jsonResult(enums.JRCodeSucc, "", data)
	} else if param == "checkchunk" {
		//检测分片是否存在
		chunk := this.GetString("chunk", "0") //分块下标
		chunkSize, _ := this.GetInt64("chunksize", 0)
		taskId := this.GetString("taskid")
		data, err := utils.CheckChunk(rootPath, fileTempPath, taskId, chunk, chunkSize)
		if err != nil {
			this.jsonResult(enums.JRCodeFailed, "", data)
		}
		this.jsonResult(enums.JRCodeSucc, "", data)
	} else if param == "mergechunks" {
		//合并文件分片
		this.mergeBlock(taskId, fileTempPath, uploadPath)
	} else {
		//fileBytes, _ := ioutil.ReadAll(file)
		//beego.Debug("fileBytes")
		//beego.Debug(fileBytes)

		////文件id   上传多个图片 id不同  WU_FILE_1、 WU_FILE_2
		//Id := this.Ctx.Request.Form.Get("id")
		//beego.Debug(Id)
		////文件类型
		//fileType := this.Ctx.Request.Form.Get("type")
		//beego.Debug(fileType) //image/png
		fileName := this.getUploadFileName()
		chunk := this.GetString("chunk", "0")             //分块下标
		saveName := fileName + "_" + taskId + "_" + chunk //文件保存名称
		this.uploadfile(uploadPath, fileTempPath, saveName)
	}
}

func (this *PictureController) getUploadFileName() string {
	//上传图片
	file, fileHead, fileErr := this.Ctx.Request.FormFile("file") //上传的文件
	//beego.Debug(file)
	//beego.Debug(fileHead)
	if file == nil {
		beego.Error(" 未找到要上传的文件!, ERROR: FormFile获取上传文件失败!")
		msg := ""
		this.jsonResult(enums.JRCodeFailed, msg, nil)
	}
	// Get a file name
	var fileName string
	if name := this.Ctx.Request.Form.Get("name"); name != "" {
		fileName = name
	} else if fileErr == nil {
		fileName = fileHead.Filename
	} else {
		fileName = utils.UniqueId() //生成一个唯一ID
	}
	return fileName
}

/**
 * 	检测分片是否存在
 *	params string tmpPath  分片存储缓存路径
 *	chunk 分块下标
 *	chunksize 分块大小
 */
func (this *PictureController) checkChunk(rootPath, tmpPath, taskId string) {
	chunk := this.GetString("chunk", "0") //分块下标
	data := make(map[string]bool)

	infoPath := tmpPath + taskId + "info"
	tempInfo, n := utils.GetJsonFileInfo(infoPath)
	if n == 0 {
		data["isExist"] = false
		this.jsonResult(enums.JRCodeFailed, "tempinfo is null", nil)
	}

	if chunkSize, ok := interface{}(this.GetString("chunksize", "0")).(int64); !ok {
		taskId := this.GetString("taskid")
		if !this.Empty(chunk) || !this.Empty(chunkSize) {
			beego.Error("分块下标或分块大小不能为空! ERROR: chunk or chunkSize is empty!")
			this.jsonResult(enums.JRCodeFailed, "invalid param", nil)
		}
		//$isExist = filesize($tmpfile) == $chunkSize;
		tempFile := rootPath + tmpPath + tempInfo.FileName + "_" + taskId + "_" + chunk + ".tmp"

		if !utils.IsFile(tempFile) || utils.GetFile(tempFile).Size() != chunkSize {
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
func (this *PictureController) mergeBlock(taskId, tmpPath, uploadPath string) {
	infoPath := tmpPath + taskId + "info"
	if utils.IsFile(infoPath) {
		//data := ; //获取文件中数据
		tempInfoFile, err := os.Open(infoPath)
		defer tempInfoFile.Close()
		if err != nil {
			this.ReturnFailedJson(err, "Failed to find file!")
		}
		tempInfoSize := utils.GetFileSize(infoPath)
		if tempInfoSize == 0 {
			beego.Error("mergeBlock: 缓存文件大小为0")
			return
		}
		data := TempInfo{}
		var info = make([]byte, tempInfoSize)
		if _, err := tempInfoFile.Read(info); err == nil {
			error := json.Unmarshal(info, &data)
			if error != nil {
				this.ReturnFailedJson(error, "Unmarshal Failed!")
			}
		} else {
			this.ReturnFailedJson(err, "Failed to find file!")
		}

		saveDir := utils.Date("Ymd", time.Now()) + string(os.PathSeparator)
		if !utils.IsDir(uploadPath + saveDir) {
			if err := os.Mkdir(uploadPath+saveDir, os.ModePerm); err != nil {
				beego.Error(errors.New("上传目录" + uploadPath + saveDir + "创建失败!"))
				this.jsonResult(enums.JRCodeFailed, "上传目录"+uploadPath+saveDir+"创建失败!", nil)
			}
		}
		//截取文件后缀名
		//index := strings.Index(data.FileName, ".")
		//ext := utils.String(data.FileName[index:])
		//创建、打开上传文件
		uploadFilePath := uploadPath + saveDir + data.FileName

		//判断文件是否已存在
		if utils.IsFile(uploadFilePath) {
			this.jsonResult(enums.JRCodeFailed, "文件已存在", nil)
		}
		uploadFile, err := os.OpenFile(uploadFilePath, os.O_CREATE|os.O_WRONLY, 0777)
		defer uploadFile.Close()
		if err != nil {
			beego.Error("mergeBlock: " + err.Error())
			this.jsonResult(enums.JRCodeFailed, "Failed to open upload file ", nil)
		}

		//锁住文件后合并缓存文件
		i := 0 //文件个数   循环分片
		var size int64 = 0
		fileSize := data.FileSize
		//优化 -> 加文件锁  ????

		//合并缓存文
		ioW := bufio.NewWriter(uploadFile) //创建新的 Writer 对象
		for size < fileSize {
			chunkFile := uploadPath + tmpPath + data.FileName + "_" + taskId + "_" + utils.String(i) + ".tmp"
			i++
			chunkFileSize := utils.GetFileSize(chunkFile)
			beego.Debug("chunkFileSize")
			beego.Debug(chunkFileSize)
			if chunkFileSize == 0 {
				utils.DeleteFile(chunkFile)
				continue
			}
			tempFile, err := os.OpenFile(chunkFile, os.O_RDONLY, 0777)
			defer tempFile.Close()
			if err != nil {
				beego.Error("mergeBlock: " + err.Error())
				break
			}
			//将缓存文件内容读取后写入上传文件中
			var buff = make([] byte, chunkFileSize)
			if _, err := tempFile.Read(buff); err != nil {
				beego.Error("mergeBlock: " + err.Error())
				break
			}
			chunkSize, err := ioW.Write(buff) //缓存块大小
			if err != nil {
				beego.Error("mergeBlock: " + err.Error())
				break
			}
			//刷新缓存buff
			ioW.Flush()
			size += int64(chunkSize)
			//合并成功, 关闭缓存文件
			if err := tempFile.Close(); err != nil {
				beego.Error("mergeBlock: " + err.Error())
				break
			}
			//删除文件
			utils.DeleteFile(chunkFile)
		}
		if size < fileSize {
			this.jsonResult(enums.JRCodeFailed, "Failed to operate the shard file "+utils.String(size)+"--"+utils.String(data.FileSize), nil)
		}

		var result = make(map[string]interface{})
		result["savePath"] = saveDir
		result["saveName"] = saveDir
		result["fileName"] = data.FileName //title
		result["uid"] = this.userId        //title
		//$data['label'] = $this->request->getQuery('label');
		//$data['audit_statu']=$this->request->getQuery('admin');
		//$data['artwork']=$this->request->getQuery('artwork');
		//$filedata = $this->_saveToDb($data);
		//$filedata['fileid'] = $data['fileid'];
		//删除记录合并信息的文件
		tempInfoFile.Close()
		utils.DeleteFile(infoPath)
		this.jsonResult(enums.JRCodeSucc, "", result)
	} else {
		beego.Error("Merge Sharding Failed!")
		this.jsonResult(enums.JRCodeFailed, "Merge Sharding Failed", nil)
	}
}

/**
 * 	上传文件    上传文件生成缓存
 */
func (this *PictureController) uploadfile(rootPath, tmpPath, saveName string) {
	beego.Debug("uploadfile 上传文件生成缓存")
	if hasFiles := this.Ctx.Request.ParseMultipartForm(32 << 20); hasFiles != nil {
		this.jsonResult(enums.JRCodeFailed, "上传文件解析失败", nil)
	}
	var fileHeads []*multipart.FileHeader
	fileHeads = this.Ctx.Request.MultipartForm.File["file"] //获取上传的文件句柄   type: array
	//上传upload类初始化
	uploadConf := make(map[string]interface{})
	uploadConf["RootPath"] = rootPath //图库根路径
	uploadConf["SavePath"] = tmpPath  //分块缓存文件存储路径
	uploadConf["AutoSub"] = false
	uploadConf["SaveName"] = saveName //缓存文件保存名称
	uploadConf["SaveExt"] = ".tmp"    //缓存文件后缀
	if upload, err := utils.NewUpload(uploadConf); err == nil {
		var infos [] map[string]interface{}
		for _, fileHead := range fileHeads {
			beego.Debug("upload temp file")
			beego.Debug(upload)
			if info, err := upload.Upload(fileHead); err != nil {
				beego.Debug(err.Error())
				this.jsonResult(enums.JRCodeFailed, fileHead.Filename+"upload failed!", nil)
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
