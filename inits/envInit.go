package inits

import "forum/pkg/globals"

func EnvInit() {
	if len(globals.Env) == 0 {
		globals.Env = "local"
	}

}
