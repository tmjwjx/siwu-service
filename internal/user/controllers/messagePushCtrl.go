package controllers

import (
	"forum/internal/user/logics"
	"github.com/gin-gonic/gin"
)

// MessagePushCtrl 向用户实时发送更新数据
func MessagePushCtrl(c *gin.Context) {

	// 设置sse响应的响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream") // 标志了该响应为事件流类型
	c.Writer.Header().Set("Cache-Control", "no-cache")         // 提示用户不要缓存响应
	c.Writer.Header().Set("Connection", "keep-alive")          // 保持连接不断开，以便持续发送事件
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // 设置跨域资源共享头，允许所有域访问该资源

	// 创建用户的消息通道 监听新消息并发送到客户端
	logics.NewMessageChan(c)

	// 监听通道中的新评论消息
	// 在外面实现
}

// LikeMessageCtrl
// @Description: 点赞消息
// @param        c *gin.Context
func LikeMessageCtrl(c *gin.Context) {

	//// 初始化需要的变量
	//db := globals.DB
	//
	//// 获取用户ID
	//userId, _ := c.Get("id")
	//fmt.Println(userId)
	//
	//// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	//if err := c.ShouldBind(&req); err != nil {
	//	// 日志记录错误信息
	//	globals.Log.Errorf("绑定req失败 err = %s", err)
	//	// 返回错误响应
	//	data := response.NewAppErr(globals.StatusBadRequest, err, nil)
	//	response.Failed(c, http.StatusBadRequest, data)
	//	return // 结束函数执行
	//}
	//fmt.Printf("%v", req)
	//
	//// 进入业务层
	//id, err := logics.ArticleCreateLogic(db, req, userId.(uint))
	//if err != nil {
	//	globals.Log.Errorf("创建文章失败 err = %s", err)
	//	data := response.NewAppErr(globals.StatusInternalServerError, err, id)
	//	response.Failed(c, http.StatusInternalServerError, data)
	//	return
	//}
	//
	//// 返回响应
	//data := response.NewAppData(globals.StatusOK, "成功", id)
	//response.Success(c, http.StatusOK, data)
}