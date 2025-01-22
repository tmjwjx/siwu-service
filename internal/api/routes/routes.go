package routes

import (
	"forum/internal/api/controlles"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

func Api(e *gin.Engine) {

	//casbinService, err := casbin.NewCasbinService(globals.DB)
	//if err != nil {
	//	fmt.Println("Api(e *gin.Engine) -> 创建 casbinService 失败, err = ", err)
	//}

	r0 := e.Group("/acl").Use(token.AuthMiddleware())
	// 获取当前api详情
	r0.GET("/api/detail", controlles.GetApiDetailsCtrl)

	// 获取所有api列表
	r0.GET("/api/list", controlles.GetAllApiCtrl)

	r := e.Group("/api").Use(token.AuthMiddleware()).Use(casbin.CasbinAuth(globals.CasbinEnforcer))
	{
		// 获取所有api分组列表
		r.GET("/groups", controlles.GetGroupListCtrl)

		// 获取所有请求方法
		r.GET("/get_methods", controlles.GetRequestMethodCtrl)

		// 删除api
		r.DELETE("/delete", controlles.DeleteApiCtrl)

		// 编辑api
		r.POST("/update", controlles.UpdateApiCtrl)

		// 添加api
		r.POST("/add", controlles.CreateApiCtrl)

		// 检索api获取列表
		r.POST("/list", controlles.SearchApiListCtrl)

	}
}
