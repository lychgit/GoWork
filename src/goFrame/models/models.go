package models

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code interface{} `json:"code"`
	Msg  string      `json:"msg"`
	Obj  interface{} `json:"obj"`
}
