package utils

import (
	"reflect"
	"strings"
	"errors"
)

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

func Empty(data interface{}) bool {
	if data == nil || data == "" {
		return true
	}
	return false
}

func String(data interface{}) string{
	return interface{}(data).(string)
}

func  InArray(data map[interface{}]interface{}, value interface{})  bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func KeyInArray(data map[interface{}]interface{}, key interface{}) bool {
	for k, _ := range data {
		if k == key {
			return true
		}
	}
	return false
}