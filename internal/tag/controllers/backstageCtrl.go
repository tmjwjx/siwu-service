package controllers

import (
	"forum/internal/tag/logics"
	"forum/internal/tag/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// AddTagCtrl 新增标签
func AddTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsAddTagReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err, status := logics.AddTagLogic(c, globals.DB, &req)

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
		d := response.NewAppData(globals.StatusOK, "新增标签成功", nil)
		response.Success(c, 200, d)
	}

}

// DeleteTagCtrl 删除标签
func DeleteTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsDelTagReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.DeleteTagLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		// 删除失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 删除成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "标签删除成功", nil)
	response.Success(c, 200, d)

}

// BatchDelTagCtrl 批量删除标签
func BatchDelTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsBatchDelTagReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.BatchDelTagLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		// 删除失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 删除成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "标签批量删除成功", nil)
	response.Success(c, 200, d)

}

// UpdateTagCtrl 更新标签
func UpdateTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsUpTagReq
	err := c.ShouldBind(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err, status := logics.UpdateTagLogic(c, globals.DB, &req)

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
		d := response.NewAppData(globals.StatusOK, "标签信息更新成功", nil)
		response.Success(c, 200, d)
	}

}

// QueryTagCtrl 查询标签
func QueryTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsQueTagReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	tagRes, err := logics.QueryTagLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		// 查询失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 查询成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "标签查询成功", tagRes)
	response.Success(c, 200, d)

}

// BatchQueryTagCtrl 批量查询标签
func BatchQueryTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.BsBatchQueTagReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	batchTagRes, err := logics.BatchQueryTagLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		// 批量查询失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 批量查询成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "标签批量查询成功", batchTagRes)
	response.Success(c, 200, d)
}
