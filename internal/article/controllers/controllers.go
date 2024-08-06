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

func SearchHandler(c *gin.Context) {
	
	// 初始化需要的变量
	db := globals.DB
	var req requests.ReqSearch // 创建一个 SearchRequest 类型的变量，用于存储请求参数
	
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
	articles, err := logics.Search(db, req)
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

func PublishHandler(c *gin.Context) {
	
	// 初始化需要的变量
	db := globals.DB
	var req requests.ReqPublish
	fmt.Println(db, req)
	
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
	err := logics.Public(db, req)
	if err != nil {
		globals.Log.Errorf("搜索文章失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		//c.JSON(500, response.StatusInternalServerErr)
		return
	}
	
	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}