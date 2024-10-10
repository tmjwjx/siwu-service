package inits

import "forum/pkg/globals"

// EnvInit 初始化环境
func EnvInit() {
	if len(globals.Env) == 0 {
		globals.Env = "local.yaml"
	}
}
