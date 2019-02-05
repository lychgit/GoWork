package utils

import (
	"github.com/astaxie/beego"
	"reflect"
	"errors"
	"strings"
)

type Upload struct {
	Mimes        map[string]interface{}      //允许上传的文件MiMe类型
	MaxSize      int32                       //上传的文件大小限制 (0-不做限制)
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
}

/**
 *	生成一个上传文件的upload结构体
 *  默认配置   利用反射给结构体赋值
 */
func NewUpload(params map[string]interface{}) (Upload, error) {
	beego.Debug("Construct")
	//初始化文件上传配置
	upload := Upload{
		MaxSize: 0,    //上传的文件大小限制 (0-不做限制)
		AutoSub: true, //自动子目录保存文件
		SubName:     make(map[int]string),             //子目录创建方式，[0]-函数名，[1]-参数，多个参数使用数组  'subName' => array('date', 'Ymd'),
		RootPath: "/upload", // BASE_DIR . '/public/uploads/', //保存根路径
		SavePath: "",        //保存路径
		SaveExt:  "",        //文件保存后缀，空则使用原后缀
		Replace:  false,     //存在同名是否覆盖
		Hash:     true,      //是否生成hash编码
		CallBack: false,     //检测文件是否存在回调，如果存在返回文件信息数组
		Driver:  "",        // 文件上传驱动
	}
	upload.SubName[0] = "date"
	upload.SubName[1] = "Ymd"


	ref := reflect.ValueOf(&upload)
	beego.Debug("ref")
	beego.Debug(ref.Kind())
	if ref.Kind() == reflect.Struct ||  ref.Kind() == reflect.Ptr{
		elem := ref.Elem()
		beego.Debug(elem)
		for i:=0; i < elem.NumField(); i++ {
			structField := elem.Type().Field(i)
			beego.Debug("structField")
			beego.Debug(structField)
			tag := structField.Tag
			beego.Debug("tag")
			beego.Debug(tag)
			structKey := tag.Get("json")
			beego.Debug("structKey")
			beego.Debug(structKey)
			if structKey == "" {
				structKey = strings.ToLower(structField.Name)
			}
			//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
			structKey = strings.Split(structKey, ",")[0]
			beego.Debug("JSONnAME:", structKey)

			if value, k := params[structKey]; k {
				beego.Debug("value")
				beego.Debug(value)
				beego.Debug(" k ")
				beego.Debug(k)
				//给结构体赋值
				//保证赋值时数据类型一致
			}

		}
		return upload, nil
	} else {
		return Upload {}, errors.New("Ptr type error")
	}
}

func Gstruct(upload Upload) interface{} {
	t := reflect.TypeOf(upload)
	if t.Kind() != reflect.Struct {
		return nil
	}
	return nil
}
