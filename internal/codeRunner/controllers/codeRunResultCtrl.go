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

// GetCodeRunResultCtrl 从codeRunner服务中获取代码运行结果
func GetCodeRunResultCtrl(c *gin.Context) {
	fmt.Println(1111111)
	// 接收数据
	req := new(requests.ExecuteResponse)
	fmt.Println(req)
	err := c.ShouldBindJSON(req)
	if err != nil {
		globals.Log.Errorf("获取失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 进入业务层
	if err := logics.NewResultSendLogic(req).Send(); err != nil {
		globals.Log.Errorf("消息发送失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	response.Success(c, http.StatusOK, nil)
}
