package response

import "forum/pkg/globals"



type AppData struct {
	Code globals.AppCode `json:"code"`
	Msg  string          `json:"msg"`
	Data interface{}     `json:"data"`
}

type AppErr struct {
	Code globals.AppCode `json:"code"`
	Err  error           `json:"err"`
	Data interface{}     `json:"data"`
}

// NewAppData 生产一个成功消息响应结构体
func NewAppData(code globals.AppCode, msg string, data interface{}) *AppData {
	return &AppData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// NewAppErr 生产一个失败消息响应结构体
func NewAppErr(code globals.AppCode, err error, data interface{}) *AppErr {
	return &AppErr{
		Code: code,
		Err:  err,
		Data: data,
	}
}
