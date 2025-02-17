package controllers

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/tag/logics"
	"forum/internal/tag/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"forum/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UpdateTagUserCount 更新标签的关注人数
func UpdateTagUserCount(c *gin.Context) {

	// 获取参数
	var tag requests.TagReq
	if err := c.ShouldBind(&tag); err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserDataResponseCtrl -> 从token中获取用户ID失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	uintValue, err := internalUtils.ChangeAnyToUint(userID)
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	// 业务处理
	fansCount, err := logics.UpdateTagUserCountLogic(uintValue, globals.DB, tag.ID)

	// 返回响应
	if err != nil {
		// 人数更新失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 人数更新成功，返回现在人数
	d := response.NewAppData(globals.StatusOK, "人数更新成功", fansCount)
	response.Success(c, 200, d)

}

// UpdateTag 刷新前端标签页
func UpdateTag(c *gin.Context) {

	// 获取参数
	userID := c.Query("user_id")
	uintValue, err := utils.ChangeStringToUint(userID)
	// 返回响应
	if err != nil {
		// 更新失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	//userID, exists := c.Get("id")
	//if !exists {
	//	// 返回错误响应
	//	e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserDataResponseCtrl -> 从token中获取用户ID失败"), nil)
	//	response.Failed(c, 500, e)
	//	return
	//}
	//
	//uintValue, err := internalUtils.ChangeAnyToUint(userID)
	//if err != nil {
	//	// 返回错误响应
	//	e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
	//	response.Failed(c, 500, e)
	//	return
	//}
	// 业务处理

	tagRes, err := logics.UpdateTagArticleCountLogic(uintValue, globals.DB) // 更新前端的标签页

	// 返回响应
	if err != nil {
		// 更新失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 更新成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "标签页数据更新成功", tagRes)
	response.Success(c, 200, d)

}

// StorageTagCtrl
// @Description:存储新用户选择的标签
// @Author wangyulong 2024-10-14 21:24:00
// @param        c *gin.Context
func StorageTagCtrl(c *gin.Context) {

	// 获取参数
	var req requests.StorageTagReq
	if err := c.ShouldBind(&req); err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 从token中获取用户id
	gid, exists := c.Get("id")
	if !exists {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("从token中获取用户id失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 将 any 转换成 uint
	uid, ok := gid.(uint)
	if !ok {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("any转换uint失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 业务处理
	err := logics.StorageTagLogic(globals.DB, &req, uid)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	d := response.NewAppData(globals.StatusOK, "人数更新成功", nil)
	response.Success(c, 200, d)

}

// GetAllTagCtrl
// @Description: 获取所有标签的id和name
// @Author wangyulong 2024-10-14 21:59:07
// @param        c *gin.Context
func GetAllTagCtrl(c *gin.Context) {

	// 逻辑处理
	res, err := logics.GetAllTagLogic(globals.DB)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	d := response.NewAppData(globals.StatusOK, "标签页数据更新成功", res)
	response.Success(c, 200, d)

}
