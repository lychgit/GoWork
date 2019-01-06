package models

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code int        `json:"code"`
	Msg  string      `json:"msg"`
	Obj  interface{} `json:"obj"`
}
