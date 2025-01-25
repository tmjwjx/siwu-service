package controllers

import (
	"forum/internal/article/logics"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ArticleBanCtrl
// @Description: 封禁文章
// @param        c *gin.Context
func ArticleBanCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	idList := requests.ArticleOperationListReq{}

	// 绑定查询参数到变量
	if err := c.BindJSON(&idList); err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 进入业务层
	err := logics.ArticleBanLocal(db, idList)
	if err != nil {
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}

// ArticleUnblockCtrl
// @Description: 解封文章
// @param        c *gin.Context
// @Author tianjiajie 2025-01-21 11:30:06
func ArticleUnblockCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	idList := requests.ArticleOperationListReq{}

	// 绑定查询参数到变量
	if err := c.BindJSON(&idList); err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 进入业务层
	err := logics.ArticleUnblockLocal(db, idList)
	if err != nil {
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}

// DeleteArticlesCtrl
// @Description: 删除文章
// @param        c *gin.Context
func DeleteArticlesCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	idList := requests.ArticleOperationListReq{}

	// 绑定查询参数到变量
	if err := c.BindJSON(&idList); err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 进入业务层
	err := logics.ArticleDeleteLocal(db, idList)
	if err != nil {
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}

// UserDeleteArticlesCtrl
// @Description: 用户 删除文章
// @param        c *gin.Context
// @Author tianjiajie 2025-01-21 09:04:12
func UserDeleteArticlesCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 绑定查询参数到变量
	articleId := c.Query("id")
	userId, ok := c.Get("id")
	if !ok {
		data := response.NewAppErr(globals.StatusInternalServerError, nil, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 进入业务层
	err := logics.UserArticleDeleteLocal(db, articleId, userId.(uint))
	if err != nil {
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}
