package controllers

import (
	"fmt"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ClickAttention 点击关注和点击取消关注
func ClickAttention(c *gin.Context) {
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var followMsg requests.ClickAttentionReq
	if err := c.ShouldBind(&followMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ClickAttention() err: %v", err), nil))
		return
	}

	follerId, exists := c.Get("id")
	if !exists {
		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("ClickAttention() err = 无法获取 id"), nil))
		return
	}
	// 类型断言
	followMsg.FollowerId = follerId.(uint)

	// 简单检验数据
	if followMsg.FollowerId == followMsg.FollowedId {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ClickAttention() : id%d不能关注%d", followMsg.FollowerId, followMsg.FollowedId), nil))
		return
	}

	// 业务逻辑
	if err := userReqContext.ClickAttention(followMsg); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("ClickAttention() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// UserRank 用户热度排行
func UserRank(c *gin.Context) {
	// 从上下文中获取 id
	str, exists := c.Get("id")
	if !exists {
		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("UserRank() err = 无法获取 id"), nil))
		return
	}
	// 类型断言
	id := str.(uint)

	// 绑定数据
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UserRank() err = 数据类型转换错误"), nil))
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UserRank() err = 数据类型转换错误"), nil))
		return
	}
	var rankMsg requests.UserRankReq
	rankMsg.Page = page
	rankMsg.Limit = limit

	// 检验数据
	if page <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UserRank() err: Page参数必须为正数"), nil))
		return
	}
	if limit <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UserRank() err: limit参数必须为正数"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	userRankRep, err := userReqContext.UserRank(id, rankMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserRank() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"user_heat_rank": userRankRep}))
}

// Attention 搜索用户关注的人
func Attention(c *gin.Context) {
	// 绑定数据
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Attention() err = 数据类型转换错误"), nil))
		return
	}
	keyword := c.Query("keyword")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Attention() err = 数据类型转换错误"), nil))
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Attention() err = 数据类型转换错误"), nil))
		return
	}

	// 检验数据
	if page <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Attention() err: Page参数必须为正数"), nil))
		return
	}
	if limit <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Attention() err: limit参数必须为正数"), nil))
		return
	}

	attentionReq := requests.AttentionReq{
		UserId:  uint(userId),
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	ids, err := userReqContext.Attention(attentionReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Attention() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"ids": ids}))
}

// GetUserArticleCtrl
// @Description: 获取用户文章
// @param        c *gin.Context
// @Author tianjiajie 2025-01-17 16:39:16
//func GetUserArticleCtrl(c *gin.Context) {
//	// 初始化需要的变量
//	db := globals.DB
//	req := requests.UserDataRequest{}
//
//	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
//	if err := c.ShouldBindQuery(&req); err != nil {
//		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
//		response.Failed(c, http.StatusBadRequest, data)
//		return
//	}
//	userId, ok := c.Get("id")
//	if ok != true {
//		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("GetUserArticleCtrl() err = 无法获取 id"), nil))
//		return
//	}
//	b := userId.(uint) == req.Id
//
//	fmt.Printf("%#v\n", req)
//	fmt.Println(userId)
//
//	// 进入业务层
//	totalData, err := logics.GetUserArticleLogic(db, req, b)
//	if err != nil {
//		globals.Log.Errorf("获取用户文章失败 err = %s", err)
//		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
//		response.Failed(c, http.StatusInternalServerError, data)
//		return
//	}
//
//	// 返回响应
//	data := response.NewAppData(globals.StatusOK, "成功", totalData)
//	response.Success(c, http.StatusOK, data)
//
//}
