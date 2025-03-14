package controllers

import (
	"fmt"
	"forum/internal/codeRunner/logics"
	"forum/internal/codeRunner/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCodeInfoCtr(c *gin.Context) {
	//接收前端的代码
	var codeInfo requests.CodeRunnerReq
	if err := c.ShouldBind(&codeInfo); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetCodeInfo() -> %v", err), nil))
		return
	}
	//调用gettoken得到token
	etcdToken := logics.CodeToken{Name: "思悟", Password: "123456"}
	token, err := etcdToken.GetEtcdToken()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetCodeInfo() -> %v", err), nil))
		return
	}
	//发送代码段
	sendCode := logics.SendCodeLogin{CodeRunnerReq: codeInfo}
	_, err = sendCode.SendCode(token)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetCodeInfo() -> %v", err), nil))
		return
	}
	//返回成功响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}
