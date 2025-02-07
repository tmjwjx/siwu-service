package controllers

import (
	"fmt"
	"forum/internal/backstage/logics"
	"forum/internal/backstage/repositories"
	"forum/internal/backstage/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"forum/pkg/response"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BsLogin
// @Description: 后台登陆
// @Author lizhuang 2024-10-21 21:21:28
// @param        c *gin.Context
func BsLogin(c *gin.Context) {
	// 绑定数据
	var bsLogicMsg requests.BackstageLoginReq
	if err := c.ShouldBind(&bsLogicMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("BsLogin() -> %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(bsLogicMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("BsLogin() : 邮箱不合法"), nil))
		return
	}
	// 检验密码是否合法
	if !internalUtils.IsValidPassword(bsLogicMsg.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("BsLogin() : 密码必须要同时包含字母、数字、特殊字符，长度在8到20位之间"), nil))
		return
	}

	// 业务逻辑
	bsManageContext := logics.NewBsManageContext(globals.DB, c)
	backstageLoginRes, err := bsManageContext.BsLogin(bsLogicMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("BsLogin() -> %v", err), nil))
		return
	}

	// 通过email查询id
	user := repositories.QueryUserByEmail(globals.DB, bsLogicMsg.Email)
	if user == nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("BsLogin() err: 不存在email为 %v 的用户", bsLogicMsg.Email), nil))
		return
	}
	// 生成token
	tok, err := token.GenerateToken(user.ID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("BsLogin() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"token": tok, "userInfo": backstageLoginRes}))
}

// BsLogout 后台登出
func BsLogout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("BsLogout() : 缺少授权标头"), nil))
		return
	}
	tokenString = tokenString[len("Bearer "):]

	// 业务逻辑
	bsManageContext := logics.NewBsManageContext(globals.DB, c)
	if err := bsManageContext.BsLogout(tokenString); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("BsLogout() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}
