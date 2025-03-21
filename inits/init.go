package inits

// init 初始化
func init() {
	// 选择环境
	EnvInit()

	// 根据环境初始化配置文件
	ConfigInit()

	// 初始化服务
	ServerInit()

	// 初始化日志文件
	LogInit()

	// 初始化 mysql
	DBInit()

	// 初始化 redis
	RedisInit()

	// 初始化表
	TableInit()

	// 初始化发送邮件配置
	SendEmailCfgInit()

	// 初始化路由配置
	RouterInit()

	// 初始化静态文件映射配置
	StaticInit()

	// 生产图片的部分路径前缀
	PartPathPrefixInit()

	// 初始化casbin配置
	CasbinInit()

	//初始化ETCD //todo 记得改回去
	InitEtcd([]string{"localhost:2379"})
}
