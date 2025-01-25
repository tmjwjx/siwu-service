package controllers

import (
	"forum/internal/article/logics"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetWorkplaceDataCtrl
// @Description: 获取工作台数据
// @param        e *gin.Engine
// @Author tianjiajie 2025-01-16 14:32:40
func GetWorkplaceDataCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	// 无需绑定

	// 进入业务层
	totalData, err := logics.GetWorkplaceDataLogic(db)
	if err != nil {
		globals.Log.Errorf("查询失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", totalData)
	response.Success(c, http.StatusOK, data)
}

// GetHotTagsCtrl
// @Description: 查询热门标签
// @param        c *gin.Context
// @Author tianjiajie 2025-01-16 14:32:19
func GetHotTagsCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	// 无需绑定

	// 进入业务层
	tags, err := logics.GetHotTagsLogic(db)
	if err != nil {
		globals.Log.Errorf("查询失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", tags)
	response.Success(c, http.StatusOK, data)
}

// GetHotArticleCtrl
// @Description: 查询前五篇热门文章数据
// @param        c *gin.Context
// @Author tianjiajie 2025-01-15 14:50:08
func GetHotArticleCtrl(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	// 无需绑定

	// 进入业务层
	// 获取文章列表
	articleList, err := logics.GetHotArticleLogic(db)
	if err != nil {
		globals.Log.Errorf("查询失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
}

// GetTwoWeeksArticleSumCtrl
// @Description: 查询近两周文章发布数量
// @Author tianjiajie 2025-01-15 09:14:26
func GetTwoWeeksArticleSumCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	// 无需绑定

	// 进入业务层
	articleSum, err := logics.GetTwoWeeksArticleSumLogic(db)
	if err != nil {
		globals.Log.Errorf("获取数据失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleSum)
	response.Success(c, http.StatusOK, data)

}
