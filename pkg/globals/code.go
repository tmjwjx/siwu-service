package globals

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200

type AppCode int

const (
	StatusOK                  AppCode = 2000 // 成功
	StatusBadRequest          AppCode = 4000 // 请求语法错误或无效参数
	StatusUnauthorized        AppCode = 4010 // 状态未经授权
	StatusForbidden           AppCode = 4230 // 状态禁止
	StatusTooManyRequests     AppCode = 4290
	StatusInternalServerError AppCode = 5000 // 服务器内部错误
)

type SSEType string

const (
	NoticeType     SSEType = "notice"
	VirtualMachine SSEType = "virtual_machine"
	CodeRunnerType SSEType = "codeRunner"
	AIType         SSEType = "ai"
)
