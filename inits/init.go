package inits

func Init() {

	// 初始化配置文件
	viperInit()

	// 初始化 MYSQL
	dbInit()

	// 初始化日志
	logInit()

}
