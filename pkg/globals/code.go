package globals

// 自定义状态码：StatusOK = 2000，区别于 http.StatusOK = 200

type AppCode int

const (
	StatusOK                  AppCode = 2000 // 成功
	StatusBadRequest          AppCode = 4000 // 请求语法错误或无效参数
	StatusUnauthorized        AppCode = 4001 // 状态未经授权
	StatusInternalServerError AppCode = 5000 // 服务器内部错误
)

// Home :图片所属单位，即属于文章图片还是用户图片或是资源表图片
type Home string

const (
	User          Home = "user"
	Article       Home = "article"
	Tag           Home = "tag"
	Category      Home = "category"
	Advertisement Home = "advertisement"
	Comment       Home = "comment"
)
