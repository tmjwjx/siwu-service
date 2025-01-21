package controllers

import (
	"fmt"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// InitUserInfoCtrl
// @Description: 初始化用户信息
// @Author wangyulong 2024-10-10 15:28:39
// @param        c *gin.Context
func InitUserInfoCtrl(c *gin.Context) {

	// 从查询参数中获取前端本页面用户的id
	authorId := c.Query("author_id")

	// 从token中获取用户id
	id, exists := c.Get("id")
	if !exists {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("InitUserInfoCtrl -> 从token中获取用户ID失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	// 逻辑处理
	req, err := logics.InitUserInfoLogic(globals.DB, fmt.Sprintf("%v", id), authorId)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "用户信息响应成功", req)
	response.Success(c, 200, d)

}

// EditSignatureCtrl
// @Description: 编辑个签
// @Author wangyulong 2024-10-14 21:18:05
// @param        c *gin.Context
func EditSignatureCtrl(c *gin.Context) {

	// 获取参数
	var req requests.EditSignatureReq
	// 绑定 JSON 数据到结构体
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// 处理绑定错误
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}
	// 从token中获取用户id
	gid, exists := c.Get("id")
	if !exists {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("EditSignatureCtrl -> 从token中获取用户id失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 将 any 转换成 uint
	uid, ok := gid.(uint)
	if !ok {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("EditSignatureCtrl -> any转换uint失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.EditSignatureLogic(globals.DB, &req, uid)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "编辑个签成功", req)
	response.Success(c, 200, d)
}
