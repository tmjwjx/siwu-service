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

// Follow 关注和取消关注
func Follow(c *gin.Context) {
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var followMsg requests.FollowReq
	if err := c.ShouldBind(&followMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Follow() err: %v", err), nil))
		return
	}

	// 简单检验数据
	if followMsg.FollowerId == followMsg.FollowedId {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Follow() : id%d不能关注%d", followMsg.FollowerId, followMsg.FollowedId), nil))
		return
	}

	// 业务逻辑
	if err := userReqContext.Follow(followMsg); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Follow() -> %v", err), nil))
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
		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("UserRank() err = 无法获取 email"), nil))
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
