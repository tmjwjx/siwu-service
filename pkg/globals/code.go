package globals

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200

type AppCode int

const (
	StatusOK                  AppCode = 2000
	StatusBadRequest          AppCode = 4000
	StatusInternalServerError AppCode = 5000
)

var CodeMsgMap = map[AppCode]string{
	StatusOK:                  "成功",
	StatusBadRequest:          "前端发送数据错误",
	StatusInternalServerError: "服务器错误",
}
