package response

import (
	"errors"
	"forum/pkg/globals"
)

// 默认信息

var (
	// StatusOkData 成功
	StatusOkData = &AppData{Code: globals.StatusOK, Msg: globals.CodeMsgMap[globals.StatusOK], Data: nil}
)

var (
	// StatusBadRequestErr 请求语法错误或无效参数
	StatusBadRequestErr = &AppErr{Code: globals.StatusBadRequest, Err: errors.New(globals.CodeMsgMap[globals.StatusBadRequest]), Data: nil}
	// StatusForbiddenErr 无权限访问资源
	StatusForbiddenErr = &AppErr{Code: globals.StatusForbidden, Err: errors.New(globals.CodeMsgMap[globals.StatusForbidden]), Data: nil}
	// StatusNotFoundErr 请求的资源不存在
	StatusNotFoundErr = &AppErr{Code: globals.StatusNotFound, Err: errors.New(globals.CodeMsgMap[globals.StatusNotFound]), Data: nil}
	// StatusInternalServerErr 服务器内部错误
	StatusInternalServerErr = &AppErr{Code: globals.StatusInternalServerError, Err: errors.New(globals.CodeMsgMap[globals.StatusInternalServerError]), Data: nil}
)
