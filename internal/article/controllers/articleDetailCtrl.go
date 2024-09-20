package controllers

import (
	"forum/internal/article/logics"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ArticleDetailCtrl(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB
	id := c.Param("id")

	// 进入业务层
	articleDetail, err := logics.ArticleDetailLogic(db, id)
	if err != nil {
		globals.Log.Errorf("加载文章详情失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		//c.JSON(500, response.StatusInternalServerErr)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleDetail)
	response.Success(c, http.StatusOK, data)

}