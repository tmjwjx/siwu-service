package response

import "forum/pkg/globals"

type AppData struct {
	Code globals.AppCode `json:"code"`
	Msg  string          `json:"msg"`
	Data interface{}     `json:"data"`
}
