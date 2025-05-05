package controllers

import (
	"encoding/json"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/virtualMachine/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateVMCallback
// @Description:
// @param        c *gin.Context
// @Author tianjiajie 2025-05-07 22:02:59
func CreateVMCallback(c *gin.Context) {
	db := globals.DB

	cb := requests.CreateVMCallback{}
	if err := c.ShouldBindJSON(&cb); err != nil {
		return
	}
	data, _ := json.Marshal(cb)

	// 通过email查找用户id
	user := internalUtils.QueryUserByEmail(db, cb.Email)
	userId := strconv.Itoa(int(user.ID))

	internalUtils.MessagePush2(string(data), userId, globals.VirtualMachineType, "")
}
