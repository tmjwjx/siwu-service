package controllers

import (
	"fmt"
	"forum/internal/message/logics"
	"forum/internal/message/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CollectionUnreadCtrl
// @Description: 未读收藏消息数量
// @param        c *gin.Context
// @Author tianjiajie 2025-02-12 22:08:00
func CollectionUnreadCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 获取用户ID
	userId, _ := c.Get("id")

	//进入业务层
	count, err := logics.CollectionUnreadCountLogic(db, userId.(uint))
	if err != nil {
		globals.Log.Errorf("收藏消息加载失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "收藏消息加载成功", count)
	response.Success(c, http.StatusOK, data)
}

// CollectionMessageCtrl
// @Description: 点赞消息
// @param        c *gin.Context
// @Author tianjiajie 2024-10-05 17:27:22
func CollectionMessageCtrl(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB
	var req requests.MessageReq

	// 获取用户ID
	userId, _ := c.Get("id")
	globals.Log.Infoln("userId:", userId)

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}
	fmt.Printf("%v", req)

	//进入业务层
	collectionList, err := logics.CollectionMessageLogic(db, req, userId.(uint))
	if err != nil {
		globals.Log.Errorf("收藏消息加载失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "收藏消息加载成功", collectionList)
	response.Success(c, http.StatusOK, data)
}
