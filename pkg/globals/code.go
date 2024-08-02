package globals

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200

type AppCode int

const (
	StatusOK                  AppCode = 2000
	StatusBadRequest          AppCode = 4000
	StatusForbidden           AppCode = 4003
	StatusNotFound            AppCode = 4004
	StatusInternalServerError AppCode = 5000
)

var CodeMsgMap = map[AppCode]string{
	StatusOK:                  "成功",
	StatusBadRequest:          "请求语法错误或无效参数",
	StatusForbidden:           "无权限访问资源",
	StatusNotFound:            "请求的资源不存在",
	StatusInternalServerError: "服务器内部错误",
}
