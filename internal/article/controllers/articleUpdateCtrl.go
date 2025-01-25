package controllers

import (
	"fmt"
	"forum/internal/article/logics"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LikeArticleCtrl
// @Description: 点赞文章
// @param        c *gin.Context
func LikeArticleCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	var req requests.ArticleLikeReq
	userId, _ := c.Get("id")
	fmt.Println("当前用户id", userId)

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	err := logics.LikeArticleLogic(db, req, userId.(uint))
	if err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("点赞失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return // 结束函数执行
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "点赞成功", nil)
	response.Success(c, http.StatusOK, data)
}

func CollectionCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	var req requests.ArticleCollectionReq
	userId, _ := c.Get("id")
	fmt.Println("当前用户id", userId)

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	err := logics.CollectionArticleLogic(db, req, userId.(uint))
	if err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("收藏失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return // 结束函数执行
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "收藏成功", nil)
	response.Success(c, http.StatusOK, data)

}