package errors

import (
	"forum/pkg/globals"
)

type AppError struct {
	Code    globals.AppCode `json:"code"`    // 自定义状态码
	Message string          `json:"message"` // 自定义错误
}

func (e *AppError) Error() string {
	return e.Message
}
