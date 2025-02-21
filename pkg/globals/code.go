package globals

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200

type AppCode int

const (
	StatusOK                  AppCode = 2000 // 成功
	StatusBadRequest          AppCode = 4000 // 请求语法错误或无效参数
	StatusUnauthorized        AppCode = 4001 // 状态未经授权
	StatusForbidden           AppCode = 4003 // 状态禁止
	StatusTooManyRequests     AppCode = 4029
	StatusInternalServerError AppCode = 5000 // 服务器内部错误
)
