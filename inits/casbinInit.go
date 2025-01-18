package inits

import (
	"fmt"
	"forum/pkg/casbin"
	"forum/pkg/globals"
)

func CasbinInit() {
	fmt.Println("*********************************************************")
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		fmt.Println("CasbinInit() -> 创建 casbinService 失败, err = ", err)
	}

	globals.CasbinEnforcer = casbinService.Enforcer

	policies, _ := casbinService.Enforcer.GetPolicy()
	fmt.Println("Casbin policies:", policies)

	roles, err := casbinService.Enforcer.GetImplicitRolesForUser("2")
	if err != nil {
		fmt.Println("Error getting roles:", err)
	} else {
		fmt.Println("Roles for user '2':", roles)
	}

	fmt.Println("*********************************************************")
}
