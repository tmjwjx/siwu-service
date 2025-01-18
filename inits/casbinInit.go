package inits

import (
	"fmt"
	"forum/pkg/casbin"
	"forum/pkg/globals"
)

func CasbinInit() {
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		fmt.Println("CasbinInit() -> 创建 casbinService 失败, err = ", err)
	}

	globals.CasbinEnforcer = casbinService.Enforcer

}
