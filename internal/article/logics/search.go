package logics

import (
	"fmt"
	"forum/internal/article/repositorys"
	"forum/internal/article/requests"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Search(c *gin.Context) {
	var req requests.SearchRequest // 创建一个 SearchRequest 类型的变量，用于存储请求参数

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid query parameters"}) // 返回 400 错误
		return                                                                     // 结束函数执行
	}
	fmt.Println(req)

	articles := repositorys.SearchArticles(c, req)

	c.JSON(200, gin.H{"data": articles, "msg": "发送成功"})
}
