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

/*
	// Ctrl模板
	// 初始化需要的变量
	db := globals.DB
	var req *requests.XXXReq

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
	articleList, err := logics.XXXLogic(db, req)
	if err != nil {
		globals.Log.Errorf("获取数据失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
*/

func GetFollowingArticleCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	req := requests.GetFollowArticleReq{}

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	userId, ok := c.Get("id")
	if !ok {
		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("无法获取 id"), nil))
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	articleList, err := logics.GetFollowingArticleLogic(db, req, userId.(uint))
	if err != nil {
		globals.Log.Errorf("获取数据失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
}

// ArticleListCtrl
// @Description: 检索获取已经发布的文章列表
// @param        c *gin.Context
func ArticleListCtrl(c *gin.Context) {

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
	articleList, err := logics.ArticleListLogic(db, req)
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

// ArticleSearchCtrl
// @Description: 搜索文章
// @param        c *gin.Context
func ArticleSearchCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	var req *requests.ArticleSearchReq // 创建一个 SearchRequest 类型的变量，用于存储请求参数

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBindQuery(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	articles, err := logics.ArticleSearchLogic(db, req)
	if err != nil {
		globals.Log.Errorf("搜索文章失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		//c.JSON(500, response.StatusInternalServerErr)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articles)
	response.Success(c, http.StatusOK, data)
}

// ArticleDetailCtrl
// @Description: 获取文章详情
// @param        c *gin.Context
func ArticleDetailCtrl(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB
	articleId := c.Query("id")

	userId, ok := c.Get("id")
	if ok == false {
		//data := response.NewAppErr(globals.StatusBadRequest, nil, nil)
		//response.Failed(c, http.StatusBadRequest, data)
		//return
		userId = uint(0)
	}

	// 进入业务层
	articleDetail, err := logics.ArticleDetailLogic(db, articleId, userId.(uint))
	if err != nil {
		globals.Log.Errorf("加载文章详情失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		//c.JSON(500, response.StatusInternalServerErr)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleDetail)
	response.Success(c, http.StatusOK, data)
}

// ArticleEditCtrl
// @Description: 返回 编辑文章界面 需要的数据
// @param        c *gin.Context
func ArticleEditCtrl(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB

	// 进入业务层
	edit, err := logics.ArticleEditLogic(db)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", edit)
	response.Success(c, http.StatusOK, data)
}

// GetArticlesByTagCtrl
// @Description: 获取标签下的文章
// @param        c *gin.Context
// @Author tianjiajie 2024-10-15 16:55:33
func GetArticlesByTagCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	var req *requests.GetArticleByTagReq
	//userId := c.Query("id")
	//kind := c.Query("kind")

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBindQuery(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	articleList, err := logics.GetArticlesByTagLogic(db, req)
	if err != nil {
		globals.Log.Errorf("获取文章失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
}

// GetUserArticleOrCollectionCtrl
// @Description: 获取用户文章或收藏列表
// @param        c *gin.Context
// @Author tianjiajie 2024-10-18 15:37:21
func GetUserArticleOrCollectionCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	var req *requests.UserArticleOrCollectionReq
	userId, exists := c.Get("id")
	if !exists {
		//response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("无法获取 id"), nil))
		//return
		userId = 0
	}
	// 类型断言
	id := userId.(uint)

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBindQuery(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	articleList, err := logics.GetUserArticleOrCollectionLogic(db, req, int(id))
	if err != nil {
		globals.Log.Errorf("获取数据失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
}
