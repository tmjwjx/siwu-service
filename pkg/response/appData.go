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
