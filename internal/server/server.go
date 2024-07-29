package server

// Run 启动路由
func Run() {

	err := Router.Run("0.0.0.0:8081")
	if err != nil {
		return
	}

}
