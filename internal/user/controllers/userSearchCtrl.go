package controllers

import (
	"fmt"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRank 用户热度排行
func UserRank(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var rankMsg requests.RankMsg
	if err := c.ShouldBind(&rankMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Follow() -> %v", err), nil))
		return
	}

	// 业务逻辑
	userResponses, err := userLogic.UserRank(rankMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserRank() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", userResponses))
}
