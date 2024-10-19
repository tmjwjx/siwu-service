package inits

import "forum/pkg/globals"

// EnvInit 初始化环境
func EnvInit() {

	// 配置项目环境 本地 local.yaml
	if len(globals.Env) == 0 {
		globals.Env = "local.yaml"
	}

}