/*
@Time : 2019/3/28 13:44 
@Author : shilinqing
@File : Base
*/
package common

import "github.com/astaxie/beego"

type BaseModel struct {
	JsonResult
	BaseQueryParam
}

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code interface{} `json:"code"`
	Message  string      `json:"message"`
	Data  interface{} `json:"data"`
}

// BaseQueryParam Base查询结构
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
}

func TableName(name string) string {
	return beego.AppConfig.String("dbprefix") + name
}