package utils

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200
const (
	StatusOK                  = 2000
	StatusBadRequest          = 4000
	StatusInternalServerError = 5000
)

var CodeMsgMap = map[int]string{
	StatusOK:                  "成功",
	StatusBadRequest:          "前端发送数据错误",
	StatusInternalServerError: "服务器错误",
}
