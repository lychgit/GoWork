package utils

import (
	"reflect"
	"strings"
	"errors"
	"github.com/astaxie/beego"
	"strconv"
)

//初始化结构体
func InitStruct(ptr interface{}, params map[string]interface{}) error {
	ref := reflect.ValueOf(ptr)
	if ref.Kind() == reflect.Ptr {
		elem := ref.Elem() //获取upload结构体中个字段的值  {map[] 0 map[] true map[0:date 1:Ymd] /upload    false true false  map[]}
		for i := 0; i < elem.NumField(); i++ {
			structField := elem.Type().Field(i) //结构体字段对应的值
			tag := structField.Tag
			structKey := tag.Get("json")
			if structKey == "" {
				//structKey = strings.ToLower(structField.Name)
				structKey = structField.Name
			}
			//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
			structKey = strings.Split(structKey, ",")[0]
			if value, ok := params[structKey]; ok {
				//给结构体赋值
				//保证赋值时数据类型一致
				//beego.Debug("set value")
				//beego.Debug("value's type")
				//beego.Debug(reflect.ValueOf(value).Type())
				//beego.Debug("structField.Name Type")
				//beego.Debug(elem.FieldByName(structField.Name).Type())
				if reflect.ValueOf(value).Type() == elem.FieldByName(structField.Name).Type() {
					//beego.Debug(structField.Name)
					//beego.Debug(reflect.ValueOf(value))
					elem.FieldByName(structField.Name).Set(reflect.ValueOf(value))
				}
			}
		}
		return nil
	} else {
		return errors.New("构造失败")
	}
}

//判断data是否为空
func Empty(data interface{}) bool {
	if data == nil || data == "" {
		return true
	}
	return false
}

//将数据的类型强制转换为string类型
func String(data interface{}) string{
	var str string
	switch data.(type) {
	case int, int8, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = strconv.Itoa(interface{}(data).(int))
		break
	case float32:
		str = strconv.FormatFloat(interface{}(data).(float64),'f',7, 32)
		break
	case float64:
		str = strconv.FormatFloat(interface{}(data).(float64),'f',15, 64)
		break
	case string:
		str = interface{}(data).(string)
		break
	default:
		str = ""
	}
	return str
}

//判断数组中是否存在某个值
func  InArray(data map[interface{}]interface{}, value interface{})  bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}
//判断数组中是否存在某个键值
func KeyInArray(data map[interface{}]interface{}, key interface{}) bool {
	for k, _ := range data {
		if k == key {
			return true
		}
	}
	return false
}
/**
 * separator	必需。规定在哪里分割字符串。
 * str	必需。要分割的字符串。
 * limit 可选。规定所返回的数组元素的数目。 可能的值：
 * 大于 0 - 返回包含最多 limit 个元素的数组
 * 小于 0 - 返回包含除了最后的 -limit 个元素以外的所有元素的数组
 * 0 - 返回包含一个元素的数组
 */
func Explode(separator string, str string)  map[interface{}]interface{}{
	attr := make(map[interface{}]interface{})
	for k,v := range strings.SplitN(str, separator, -1) {
		attr[k] = v
	}
	return attr
}
//结构体数据 debug打印输出
func DebugStruct(ptr interface{}) {
	beego.Debug("DebugStruct")
	ref := reflect.ValueOf(ptr)
	if ref.Kind() == reflect.Ptr {
		elem := ref.Elem() //获取upload结构体中个字段的值  {map[] 0 map[] true map[0:date 1:Ymd] /upload    false true false  map[]}
		for i := 0; i < elem.NumField(); i++ {
			structField := elem.Type().Field(i) //结构体字段对应的值
			tag := structField.Tag
			structKey := tag.Get("json")
			if structKey == "" {
				//structKey = strings.ToLower(structField.Name)
				structKey = structField.Name
			}
			//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
			structKey = strings.Split(structKey, ",")[0]
			beego.Debug(structKey)   //结构体字段名
			beego.Debug(elem.FieldByName(structField.Name).String()) //结构体字段值

			//if value, ok := params[structKey]; ok {
			//	//给结构体赋值
			//	//保证赋值时数据类型一致
			//	//beego.Debug("set value")
			//	//beego.Debug("value's type")
			//	//beego.Debug(reflect.ValueOf(value).Type())
			//	//beego.Debug("structField.Name Type")
			//	//beego.Debug(elem.FieldByName(structField.Name).Type())
			//	if reflect.ValueOf(value).Type() == elem.FieldByName(structField.Name).Type() {
			//		//beego.Debug(structField.Name)
			//		//beego.Debug(reflect.ValueOf(value))
			//		elem.FieldByName(structField.Name).Set(reflect.ValueOf(value))
			//	}
			//}
		}
	}
}