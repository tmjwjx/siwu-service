package globals

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200

type AppCode int
type Kind int

const (
	StatusOK                  AppCode = 2000 // 成功
	StatusBadRequest          AppCode = 4000 // 请求语法错误或无效参数
	StatusInternalServerError AppCode = 5000 // 服务器内部错误
)

// 该常量用来在存储图片路径时，辨别存文章表还是用户表
const (
	User    Kind = 1
	Article Kind = 2
)
