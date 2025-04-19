package controllers

import (
	"fmt"
	"forum/internal/AI/logics"
	"forum/internal/AI/requests"
	token2 "forum/pkg/AI/token"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveArticleIDCtrl(c *gin.Context) {
	//接收前端的代码
	var req requests.SaveArticleIDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("SaveArticleIDCtrl -> %v", err), nil))
		return
	}

	// 生成token
	etcdToken := token2.CodeToken{
		GenerateTokenKey: "123456",
	}
	token, err := etcdToken.GetToken()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("SaveArticleIDCtrl -> %v", err), nil))
		return
	}

	s := logics.SaveArticleIDLogic{
		Token: token,
	}

	res, err := s.SaveArticleIDLogic(&req)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("SaveArticleIDCtrl -> %v", err), nil))
		return
	}

	//返回成功响应
	data := response.NewAppData(globals.StatusOK, "成功", res)
	response.Success(c, http.StatusOK, data)

}
