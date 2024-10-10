package inits

func init() {
	// 初始化环境
	EnvInit()

	// 优先初始化配置文件（给mysql，redis赋上配置信息）
	ConfigInit()

	// 初始化 mysql
	DBInit()

	// 初始化 redis
	// RedisInit()

	// 初始化表
	TableInit()

	// 初始化发送邮件配置
	SendEmailCfgInit()

	// 初始化日志文件
	LogInit("logs", "forum")

	// 初始化路由配置
	RouterInit()

}