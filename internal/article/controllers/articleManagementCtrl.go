package controllers

import (
	"forum/internal/article/logics"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetArticleList
// @Description: 检索获取已经发布的文章列表
// @param        c *gin.Context
func GetArticleList(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB
	var req *requests.ArticleListReq

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}
	globals.Log.Info("%v", req)

	// 进入业务层
	articleList, err := logics.GetArticleList(db, req)
	if err != nil {
		globals.Log.Errorf("获取文章列表失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
}