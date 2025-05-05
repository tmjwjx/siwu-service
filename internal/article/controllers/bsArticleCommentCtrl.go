package controllers

import (
	"fmt"
	"forum/internal/article/logics"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// BatchReviewCtrl 批量审核
func BatchReviewCtrl(c *gin.Context) {

	//获取参数
	var req requests.BatchReviewReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("BatchReviewCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	res, err := logics.BatchReviewLogic(globals.DB, &req)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "评论批量审核成功", res)
	response.Success(c, 200, d)

}

// ShowCommentsListCtrl 展示评论列表(获取评论列表)
func ShowCommentsListCtrl(c *gin.Context) {

	//获取参数
	var req requests.CommentsListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ShowCommentsListCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	commentsListRes, err := logics.ShowCommentsListLogic(globals.DB, &req)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "评论列表响应成功", commentsListRes)
	response.Success(c, 200, d)

}

// AddCommentCtrl 添加评论
func AddCommentCtrl(c *gin.Context) {

	//获取参数
	var req requests.AddCommentReq
	err := c.ShouldBind(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddCommentCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err, status := logics.AddCommentLogic(c, globals.DB, &req)

	// 返回响应
	if err != nil {
		var state globals.AppCode
		if status == 400 {
			state = globals.StatusBadRequest
		} else if status == 500 {
			state = globals.StatusInternalServerError
		}
		// 返回错误响应
		e := response.NewAppErr(state, err, nil)
		response.Failed(c, status, e)
	} else {
		d := response.NewAppData(globals.StatusOK, "添加评论成功", nil)
		response.Success(c, 200, d)
	}

}

// BsDeleteCommentCtrl 删除评论
func BsDeleteCommentCtrl(c *gin.Context) {

	// 获取参数
	var req requests.DelCommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("BsDeleteCommentCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.BsDeleteCommentLogic(globals.DB, &req)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "评论删除成功", nil)
	response.Success(c, 200, d)

}

// BatchDelTagCtrl 批量删除评论
func BatchDelTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsBatchDelCommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.BatchDelCommentLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		// 删除失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 删除成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "评论批量删除成功", nil)
	response.Success(c, 200, d)

}

// UpdateCommentCtrl 更新评论
func UpdateCommentCtrl(c *gin.Context) {

	//获取参数
	var req requests.UpdateCommentReq
	err := c.ShouldBind(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateCommentCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err, status := logics.UpdateCommentLogic(c, globals.DB, &req)

	// 返回响应
	if err != nil {
		var state globals.AppCode
		if status == 400 {
			state = globals.StatusBadRequest
		} else if status == 500 {
			state = globals.StatusInternalServerError
		}
		// 返回错误响应
		e := response.NewAppErr(state, err, nil)
		response.Failed(c, status, e)
	} else {
		d := response.NewAppData(globals.StatusOK, "更新评论成功", nil)
		response.Success(c, 200, d)
	}

}

// QueryCommentCtrl 查询某个用户的全部评论
func QueryCommentCtrl(c *gin.Context) {

	//获取参数
	var req requests.QueryCommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ShowCommentsListCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	queryCommentRes, err := logics.QueryCommentLogic(globals.DB, &req)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "查询某个用户的全部评论响应成功", queryCommentRes)
	response.Success(c, 200, d)

}
