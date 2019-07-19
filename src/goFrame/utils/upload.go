package utils

import (
	"github.com/astaxie/beego"
	"mime/multipart"
	"errors"
	"os"
	"path"
		"io"
	"strings"
	"encoding/json"
	"bufio"
	"tests/utils"
)

type Upload struct {
	Mimes        map[string]interface{}      //允许上传的文件MiMe类型
	MaxSize      int64                         //上传的文件大小限制 (0-不做限制)
	Exts         map[string]string           //允许上传的文件后缀
	AutoSub      bool                        //自动子目录保存文件
	SubName      map[int]string              //子目录创建方式，[0]-函数名，[1]-参数，多个参数使用数组  'subName' => array('date', 'Ymd'),
	RootPath     string                      // BASE_DIR . '/public/uploads/', //保存根路径
	SavePath     string                      //保存路径
	SaveName     string                      //上传文件命名规则，[0]-函数名，[1]-参数，多个参数使用数组
	SaveExt      string                      //文件保存后缀，空则使用原后缀
	Replace      bool                        //存在同名是否覆盖
	Hash         bool                        //是否生成hash编码
	CallBack     bool                        //检测文件是否存在回调，如果存在返回文件信息数组
	Driver       string                      // 文件上传驱动
	DriverConfig map[interface{}]interface{} // 上传驱动配置
	error        error
	FileType    string       //文件类型  png、jpg ...
}

/**
 *	生成一个上传文件的upload结构体
 *  默认配置   利用反射给结构体赋值
 */
func NewUpload(params map[string]interface{}) (Upload, error) {
	//初始化文件上传配置
	upload := Upload{
		MaxSize:  0,                    //上传的文件大小限制 (0-不做限制)
		AutoSub:  true,                 //自动子目录保存文件
		SubName:  make(map[int]string), //子目录创建方式，[0]-函数名，[1]-参数，多个参数使用数组  'subName' => array('date', 'Ymd'),
		RootPath: "/upload",            // BASE_DIR . '/public/uploads/', //保存根路径
		SavePath: "",                   //保存路径
		SaveExt:  "",                   //文件保存后缀，空则使用原后缀
		Replace:  false,                //存在同名是否覆盖
		Hash:     true,                 //是否生成hash编码
		CallBack: false,                //检测文件是否存在回调，如果存在返回文件信息数组
		Driver:   "",                   // 文件上传驱动
		FileType: "",
	}
	upload.SubName[0] = "date"
	upload.SubName[1] = "Ymd"
	upload.FileType = ".gif,.jpg,.jpeg,.bmp,.png,.swf,.tmp"

	if len(params) > 0 {
		//根据传入的参数初始化upload结构体
		if err := InitStruct(&upload, params); err == nil {
			return upload, nil
		} else {
			return upload, err
		}
	}
	return upload, nil
}

/**
 * 创建目录
 * @param  string $savepath 要创建的目录
 * @return boolean          创建状态，true-成功，false-失败
 */
func (u *Upload) Mkdir(savePath string) bool {
	dir := u.RootPath + savePath
	if !IsDir(dir) {
		if err := os.Mkdir(savePath, os.ModePerm); err != nil {
			u.error = errors.New("上传目录" + savePath + "创建失败!")
			return false
		}
	}
	return true
}

/**
 * 获取最后一次上传错误信息
 * @return string 错误信息
 */
func (u *Upload) GetError() string {
	return u.error.Error()
}

/**
 * 保存指定文件
 * @param  array   $file    保存的文件信息
 * @param  boolean $replace 同名文件是否覆盖
 * @return boolean          保存状态，true-成功，false-失败
 */
func (u *Upload) saveFile(file *multipart.FileHeader, fileInfo map[string]interface{}, replace bool) bool {
	temFile, err := file.Open()
	defer temFile.Close()
	if err != nil {
		u.error = errors.New("上传文件获取失败")
		return false
	} else {
		fileName := u.RootPath + String(fileInfo["savePath"]) + String(fileInfo["saveName"])
		beego.Debug("缓存文件保存名称")
		beego.Debug(fileName)
		if !replace && IsFile(fileName) {
			/* 不覆盖同名文件 */
			u.error = errors.New("存在同名文件" + String(fileInfo["saveName"]))
			return false
		} else if replace {
			/* 覆盖同名文件 */
			beego.Debug("覆盖")

		} else {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
			defer f.Close()
			if err != nil {
				u.error = errors.New("create file error:" + err.Error())
				return false
			}
			if _, err := io.Copy(f, temFile); err != nil {
				u.error = errors.New("文件上传失败! 文件名：" + file.Filename)
				return false
			}
		}
		return true
	}
}

/**
 * 检查文件大小是否合法
 * @param integer $size 数据
 */
func (u *Upload) CheckSize(size int64) bool {
	return size <= u.MaxSize
}

/**
 * 检测上传根目录
 * @return boolean true-检测通过，false-检测失败
 */
func (u *Upload) CheckRootPath() bool {
	if !(IsDir(u.RootPath) && IsWritable(u.RootPath)) {
		u.error = errors.New("上传根目录不存在！请尝试手动创建" + u.RootPath)
		return false
	} else {
		return true
	}
}

/**
 * 检测上传目录
 * @param  string $savepath 上传目录
 * @return boolean          检测结果，true-通过，false-失败
 */
func (u *Upload) CheckSavePath() bool {
	if !IsExist(u.SavePath) {
		if err := os.Mkdir(u.SavePath, os.ModePerm); err != nil {
			u.error = errors.New("上传目录" + u.SavePath + "创建失败!")
			return false
		}
	} else if !IsWritable(u.RootPath + u.SavePath) {
		u.error = errors.New("上传根目录" + u.SavePath + "不可写")
		return false
	}
	return true
}

/**
 * 检查上传的文件MIME类型是否合法
 * @param string $mime 数据
 */
//private function checkMime($mime)
//{
//return empty($this->config['mimes']) ? true : in_array(strtolower($mime), $this->mimes);
//}

/**
 * 检查上传的文件后缀是否合法
 * @param string $ext 后缀
 */
 func (u * Upload) checkExt(ext string) bool {
 	if Empty(u.FileType) {
		return true
	} else if InArray(Explode(",",u.FileType) ,strings.ToLower(ext)) {
		return true
	} else {
		return false
	}
 }


/**
 * 根据指定的规则获取文件或目录名称
 * @param  array  $rule     规则
 * @param  string $filename 原文件名
 * @return string           文件或目录名称
 */
 func (u *Upload) getName(saveName, fileName string) string{
	return saveName + fileName
 }
//private function getName($rule, $filename)
//{
//$name = '';
//if (is_array($rule)) { //数组规则
//$func = $rule[0];
//$param = (array)$rule[1];
//foreach ($param as &$value) {
//$value = str_replace('__FILE__', $filename, $value);
//}
//if($func=='uuid'){
//$random=new \Phalcon\Security\Random();
//$func=[$random,'uuid'];
//}
//$name = call_user_func_array($func, $param);
//} elseif (is_string($rule)) { //字符串规则
//if (function_exists($rule)) {
//$name = call_user_func($rule);
//} else {
//$name = $rule;
//}
//}
//return $name;
//}

/**
 * 根据上传文件命名规则取得保存文件名
 * @param string $file 文件信息
 */
 func (u *Upload) getSaveName(file *multipart.FileHeader) string {
 	rule := u.SaveName
 	var saveName string
 	if Empty(rule) {
		saveName = file.Filename
	} else {
		saveName = u.getName(rule,"")
		if Empty(saveName) {
			u.error = errors.New("File name rule error!")
			return ""
		}
	}
 	if Empty(u.SaveExt) {
		return saveName + path.Ext(file.Filename)
	} else {
		return saveName + u.SaveExt
	}
 }

/**
 * 获取子目录的名称
 * @param array $file  上传的文件信息
 */
 func (u *Upload) getSubPath(fileName string) string {
 	var subPath string
 	subPath = ""
 	rule := u.SavePath
 	if u.AutoSub && rule != "" {
 		subPath = u.getName(rule, fileName) + "/"
	}
 	if Empty(subPath) && u.Mkdir(u.SavePath + subPath) {
		return ""
	}
 	return subPath
 }

func (u *Upload) getHash(temName string) string{
	return  Md5(temName)
}

func (u *Upload) upload(file *multipart.FileHeader) (map[string]interface{},bool) {
	/* 检测上传根目录 */
	if !u.CheckRootPath() {
		return nil, false
	}

	/* 检查上传目录 */
	if !u.CheckSavePath() {
		return nil, false
	}

	var info = make(map[string]interface{})
	info["key"] = ""
	info["type"] = "" //文件类型
	info["ext"] = path.Ext(file.Filename) //文件扩展
	info["name"] = file.Filename //文件名称
	info["temName"] = u.SaveName + "_" + file.Filename //文件生成的缓存文件名称
	info["size"] = file.Size //文件大小


	/* 无效上传 */
	if Empty(file.Filename) {
		u.error = errors.New("无效上传!")
		return nil, false
	}

	/* 检查文件大小 */
	if (u.CheckSize(file.Size)) {
		u.error = errors.New("Upload file size inconsistent!")
		return nil, false
	}

	/* 检查是否合法上传 */
	//if (!$file->isUploadedFile()) {
	//	throw new \Exception('Illegal file upload!');
	//}

	/* 检查文件Mime类型 */
	//if (!$this->checkMime($info['type'])) {
	//throw new \Exception('Upload file MIME type not allowed!');
	//}

	/* 检查文件后缀 */
	if !u.checkExt(String(info["ext"])) {
		u.error = errors.New("File suffix name error!'")
		return nil, false
	}

	/* 哈希验证 */
	if u.Hash {
		info["md5"] = u.getHash(u.SaveName)
	}

	/* 调用回调函数检测文件是否存在 */
	//if($this->callback){
	//	$data = call_user_func($this->callback, $info);
	//	if ($this->callback && $data) {
	//		if (file_exists('.' . $data['path'])) {
	//		return $data;
	//		} elseif ($this->removeTrash) {
	//			call_user_func($this->removeTrash, $data);//删除垃圾据
	//		}
	//	}
	//}

	/* 生成保存文件名 */
	info["saveName"] = u.getSaveName(file);

	/* 检测并创建子目录 */
	info["savePath"] = u.SavePath + u.getSubPath(String(info["name"]))
	//
	beego.Debug("info save")
	beego.Debug(info)
	/* 保存文件 并记录保存成功的文件 */
	if  u.saveFile(file, info, u.Replace) {
		return info,true
	} else {
		return nil, false
	}
}

func (u *Upload) Upload(file *multipart.FileHeader) (map[string]interface{},error) {
	if info, k := u.upload(file); !k {
		return nil, u.error
	} else {
		return info, nil
	}
}


/**
 *  保存文件上传的任务信息
 * tmpPath  存储缓存路径
 * tmpPath 缓存信息文件路径
 * fileId  WebUpload插件 queue队列中文件的id
 * fileName 上传的文件名称
 * label
 * chunkSize 上传的分块大小
 * fileSize 上传的文件大小
 * userId 用户id
 */
func SaveTaskInfo(tmpPath, fileHash, fileId, fileName, label string, chunkSize, fileSize, userId int64) (saveInfo map[string]interface{},error error){
	var data = make(map[string]interface{});
	pathSeparator := string(os.PathSeparator)
	taskId := utils.Md5(fileHash + "_" + string(userId))
	infoPath := tmpPath + pathSeparator + taskId + "info"
	//判断分片任务信息缓存文件是否存在
	if v := utils.IsFile(infoPath); !v {
		var info  TempInfo
		info.ChunkSize = chunkSize
		info.FileHash = fileHash
		info.FileId = fileId
		info.FileName = fileName
		info.FileSize = fileSize
		info.Label = label
		temp, err := json.Marshal(info)
		if err != nil {
			beego.Error(err.Error())
		}
		//file_put_contents($infoPath, serialize($data)); //将任务信息写入infoPath目录下保存
		infoFile, err := os.OpenFile(infoPath, os.O_CREATE|os.O_WRONLY, 0644)
		defer infoFile.Close()
		if err != nil {
			//分片任务信息缓存文件创建失败
			beego.Error("create file error:", err.Error())

		}
		ioW := bufio.NewWriter(infoFile) //创建新的 Writer 对象
		_, error := ioW.WriteString(string(temp))
		if error != nil {
			beego.Error("write error", error.Error())
			data["message"] = error.Error()
			return data, error
		}
		ioW.Flush()
	}
	data["taskid"] = taskId;
	return data, nil
}

/**
 * 	检测分片是否存在
 *	params string tmpPath  分片存储缓存路径
 *	chunk 分块下标
 *	chunksize 分块大小
 */
func CheckChunk(rootPath, tmpPath, taskId, chunk string, chunkSize int64) (checkInfo map[string]interface{}, error error) {
	data := make(map[string]interface{})
	infoPath := tmpPath + taskId + "info"
	tempInfo, n := GetJsonFileInfo(infoPath)
	if n == 0 {
		data["isExist"] = false
		return nil, errors.New("tempinfo is null")
	}

	if !Empty(chunk) || !Empty(chunkSize){
		err := errors.New("ERROR: chunk or chunkSize is empty! 分块下标或分块大小不能为空!")
		data["message"] = err.Error()
		return data, err
	}
	//$isExist = filesize($tmpfile) == $chunkSize;
	tempFile := rootPath + tmpPath + tempInfo.FileName + "_" + taskId + "_"  + chunk + ".tmp"

	if !utils.IsFile(tempFile) || utils.GetFile(tempFile).Size() != chunkSize {
		data["isExist"] = false
	} else {
		data["isExist"] = true
	}
	return data, nil
}

//
///**
// * 	上传文件    上传文件生成缓存
// */
//func UploadTmpeFile(rootPath, tmpPath, saveName string)  {
//	beego.Debug("uploadfile 上传文件生成缓存")
//	if hasFiles := this.Ctx.Request.ParseMultipartForm(32 << 20); hasFiles != nil {
//		this.jsonResult(enums.JRCodeFailed, "上传文件解析失败", nil)
//	}
//	var fileHeads []*multipart.FileHeader
//	fileHeads = this.Ctx.Request.MultipartForm.File["file"]  //获取上传的文件句柄   type: array
//	//上传upload类初始化
//	uploadConf := make(map[string]interface{})
//	uploadConf["RootPath"] = rootPath //图库根路径
//	uploadConf["SavePath"] = tmpPath //分块缓存文件存储路径
//	uploadConf["AutoSub"] = false
//	uploadConf["SaveName"] = saveName //缓存文件保存名称
//	uploadConf["SaveExt"] = ".tmp" //缓存文件后缀
//	if upload, err := utils.NewUpload(uploadConf); err == nil {
//		var infos [] map[string]interface{}
//		for _, fileHead := range fileHeads {
//			beego.Debug("upload temp file")
//			beego.Debug(upload)
//			if info, err := upload.Upload(fileHead); err != nil {
//				beego.Debug(err.Error())
//				this.jsonResult(enums.JRCodeFailed, fileHead.Filename + "upload failed!", nil)
//			} else {
//				beego.Debug(info)
//				infos = append(infos, info)
//			}
//		}
//		this.jsonResult(enums.JRCodeSucc, "upload success!", infos)
//	} else {
//		beego.Error("文件上传类实例化失败！ ERROR: file upload class create error!")
//		this.jsonResult(enums.JRCodeFailed, "upload create error!", nil)
//	}
//}