package routes

import (
	"forum/internal/api/controlles"
	"github.com/gin-gonic/gin"
)

func Api(e *gin.Engine) {

	// 获取当前api详情
	e.GET("/acl/api/detail", controlles.GetApiDetailsCtrl)

	// 获取所有api列表
	e.GET("/acl/api/list", controlles.GetAllApiCtrl)

	r := e.Group("/api")
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
