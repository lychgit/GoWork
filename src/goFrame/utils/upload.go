package utils

import (
	"github.com/astaxie/beego"
	"mime/multipart"
	"errors"
	"os"
	"path"
		"io"
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
	}
	upload.SubName[0] = "date"
	upload.SubName[1] = "Ymd"

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
func (u *Upload) Save(file *multipart.FileHeader, fileInfo map[string]interface{}, replace bool) bool {
	temFile, err := file.Open()
	defer temFile.Close()
	if err != nil {
		u.error = errors.New("上传文件获取失败")
		return false
	} else {
		fileName := u.RootPath + String(fileInfo["savePath"]) + String(fileInfo["saveName"])
		if !replace && IsFile(fileName) {
			/* 不覆盖同名文件 */
			u.error = errors.New("存在同名文件" + String(fileInfo["saveName"]))
			return false
		} else if replace {
			/* 覆盖同名文件 */
			beego.Debug("覆盖")

		} else {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				u.error = errors.New("create file error:" + err.Error())
				return false
			}
			defer f.Close()
			if _, err := io.Copy(f, temFile); err != nil {
				u.error = errors.New("文件上传失败! 文件名：" + file.Filename)
				return false
			}
		}
		return true
	}
}

func (u *Upload) Upload(file *multipart.FileHeader) (map[string]interface{},error) {
	if info, k := u.UploadFile(file); !k {
		return nil, u.error
	} else {
		return info, nil
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
 func (u * Upload) checkExt(ext string) {
 	if Empty(u.SaveExt) {

	}
	 return
 }
//private function checkExt($ext)
//{
//return empty($this->config['exts']) ? true : in_array(strtolower($ext), $this->exts);
//}


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
 	var rule string
 	rule = u.SaveName
 	var saveName string
 	if Empty(rule) {
		saveName = file.Filename
	} else {
		saveName = u.getName(rule, file.Filename)
		if Empty(saveName) {
			u.error = errors.New("File name rule error!")
			return ""
		}
	}
 	var ext string
 	if Empty(u.SaveExt) {
 		ext = path.Ext(file.Filename)
	} else {
		ext = u.SaveExt
	}
	return saveName + ext
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

func (u *Upload) UploadFile(file *multipart.FileHeader) (map[string]interface{},bool) {
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
	//if (!$this->checkExt($info['ext'])) {
	//throw new \Exception('Upload file suffixes are not allowed');
	//}

	/* 哈希验证 */
	if u.Hash {
		info["md5"] = u.GetHash(u.SaveName)
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
	/* 对图像文件进行严格检测 */
	//$ext = strtolower($info['ext']);
	//if (!$_POST['chunks'] && in_array($ext, array('gif', 'jpg', 'jpeg', 'bmp', 'png', 'swf'))) {
	//$imginfo = getimagesize($file->getTempName());
	//if (empty($imginfo) || ($ext == 'gif' && empty($imginfo['bits']))) {
	//throw new \Exception('Illegal image file!');
	//}
	//}
	//
	beego.Debug("info save")
	beego.Debug(info)
	/* 保存文件 并记录保存成功的文件 */
	if  u.Save(file, info, u.Replace) {
		return info,true
	} else {
		return nil, false
	}
}

func (u *Upload) GetHash(temName string) string{
	return  Md5(temName)
}