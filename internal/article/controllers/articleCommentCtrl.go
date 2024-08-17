package controllers

import (
	"fmt"
	"forum/internal/article/logics"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// InsertCommentCtrl 将评论存入数据库中
func InsertCommentCtrl(c *gin.Context) {
	// 获取参数
	var articleCommentReq requests.ArticleCommentReq
	err := c.ShouldBindJSON(&articleCommentReq)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("CreateCommentCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.InsertCommentLogic(&articleCommentReq, globals.DB)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "评论存入数据库成功", nil)
	response.Success(c, 200, d)
}

/*// GetCommentByArticleCtrl 普通的返回评论
func GetCommentByArticleCtrl(c *gin.Context) {

	// 获取参数
	articleID := c.Param("postID")

	// 逻辑处理
	topLevelComments, err := logics.GetCommentByArticleLogic(articleID)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	}
	d := response.NewAppData(globals.StatusOK, "评论存入数据库成功", topLevelComments)
	response.Success(c, 200, d)
}*/

// GetTopLevelCommentsCtrl 返回顶级评论
func GetTopLevelCommentsCtrl(c *gin.Context) {
	// 获取参数
	var req requests.TopCommentsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetTopLevelCommentsCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	topCommentsRes, err := logics.GetTopLevelCommentsLogic(&req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "顶级评论响应成功", topCommentsRes)
	response.Success(c, 200, d)
}

// GetRepliesRep2Ctrl 返回评论回复
func GetRepliesRep2Ctrl(c *gin.Context) {
	// 获取参数
	var req requests.RepliesReq2
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetRepliesRep2Ctrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	repliesRes, err := logics.GetRepliesRep2Logic(&req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "顶级评论响应成功", repliesRes)
	response.Success(c, 200, d)
}

// DeleteCommentCtrl 删除评论
func DeleteCommentCtrl(c *gin.Context) {

	// 获取参数
	var req requests.DelComment
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DeleteCommentCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.DeleteCommentLogic(&req, globals.DB)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "评论删除成功", nil)
	response.Success(c, 200, d)

}

// UpdatePraiseCountCtrl 更新点赞的数量
func UpdatePraiseCountCtrl(c *gin.Context) {

	// 获取参数
	var req requests.PraiseCount
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdatePraiseCountCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.UpdatePraiseCountLogic(&req, globals.DB)

	//返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "点赞数量更新成功", nil)
	response.Success(c, 200, d)

}
