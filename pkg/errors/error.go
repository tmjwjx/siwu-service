package errors

import "forum/pkg/globals"

// 常见错误
var (
	ErrUserNotFound = &AppError{Code: globals.StatusInternalServerError, Message: "用户没有被发现"}
	ErrHttpBad      = &AppError{Code: globals.StatusBadRequest, Message: globals.CodeMsgMap[globals.StatusBadRequest]}
)
