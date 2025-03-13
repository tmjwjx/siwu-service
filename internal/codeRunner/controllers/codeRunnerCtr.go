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
	var codeInfo requests.CodeRunnerReq
	if err := c.ShouldBind(&codeInfo); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetCodeInfo() -> %v", err), nil))
		return
	}
	etcdToken := logics.CodeToken{Name: "思悟", Password: "123456"}
	token, err := etcdToken.GetEtcdToken()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetCodeInfo() -> %v", err), nil))
		return
	}
	sendCode := logics.SendCodeLogin{CodeRunnerReq: codeInfo}
	_, err = sendCode.SendCode(token)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetCodeInfo() -> %v", err), nil))
		return
	}
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}
