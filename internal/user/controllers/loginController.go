package controllers

import (
	"encoding/json"
	"fmt"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 注册、登陆、验证码

// Register 注册
func Register(c *gin.Context) {
	// 获取数据包
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": globals.StatusInternalServerError, "data": "", "errors": globals.CodeMsgMap[globals.StatusInternalServerError]})
		return
	}
	fmt.Println("Register 接收到的消息为:")
	fmt.Println(string(data))

	// 将数据包反序列化
	var registerMsg requests.RegisterMsg
	err = json.Unmarshal(data, &registerMsg)
	if err != nil {
		// 反序列化失败，前端发送的数据有问题
		c.JSON(http.StatusBadRequest, gin.H{"code": globals.StatusBadRequest, "data": "", "errors": globals.CodeMsgMap[globals.StatusBadRequest]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": globals.StatusOK, "data": gin.H{}, "msg": globals.CodeMsgMap[globals.StatusOK]})
}
