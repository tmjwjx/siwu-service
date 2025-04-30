package controllers

import (
	"encoding/json"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/virtualMachine/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"strconv"
)

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

// DestroyVMCallback 销毁虚拟机毁掉
func DestroyVMCallback(c *gin.Context) {
	db := globals.DB

	dr := requests.DestroyResp{}
	if err := c.ShouldBindJSON(&dr); err != nil {
		globals.Log.Panicf("绑定数据失败: %v\n", err)
		return
	}

	// 通过email查找用户id
	user := internalUtils.QueryUserByEmail(db, dr.Email)
	if user != nil {
		userId := strconv.Itoa(int(user.ID))
		// 调用sse
		internalUtils.MessagePush2(fmt.Sprintf("已销毁Email为%v用户的虚拟机", dr.Email), userId, globals.VirtualMachineType)
	} else {
		globals.Log.Panicf("无法获取email为%v的用户id\n", dr.Email)
	}
}

// NoticeVMCallback 虚拟机信息通知
func NoticeVMCallback(c *gin.Context) {
	db := globals.DB

	noticeResp := requests.NoticeResp{}
	if err := c.ShouldBindJSON(&noticeResp); err != nil {
		globals.Log.Panicf("绑定数据失败: %v\n", err)
		return
	}

	// 通过email查找用户id
	user := internalUtils.QueryUserByEmail(db, noticeResp.Email)
	if user != nil {
		userId := strconv.Itoa(int(user.ID))
		// 调用sse
		internalUtils.MessagePush2(noticeResp.Data, userId, globals.VirtualMachineType)
	} else {
		globals.Log.Panicf("无法获取email为%v的用户id\n", noticeResp.Email)
	}
}
