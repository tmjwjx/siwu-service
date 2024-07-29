package logic

import (
	"fmt"
	"forum/internal/article/repository"
	"forum/internal/article/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Search(c *gin.Context) {
	var req request.SearchRequest // 创建一个 SearchRequest 类型的变量，用于存储请求参数

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"}) // 返回 400 错误
		return                                                                    // 结束函数执行
	}
	fmt.Println(req)

	articles := repository.SearchArticles(c, req)

	//articles = append(articles, models.Article{
	//	AuthorID:                2,
	//	ArticleTitle:            "searchApple",
	//	ArticleImage:            "",
	//	ArticleLink:             "",
	//	ArticleLikesCount:       0,
	//	ArticleCollectionsCount: 0,
	//	ArticleCommentsCount:    0,
	//	ArticleViewsCount:       0,
	//	ArticleHeat:             0,
	//	ArticleAttributes:       "",
	//	ArticleCategoryID:       0,
	//	ArticleTagID:            0,
	//	ArticleSummary:          "",
	//	ArticleContent:          "",
	//	UpdatedAt:               nil,
	//	PublishedAt:             nil,
	//	Category:                models.Category{},
	//})

	// 返回查询到的产品列表，状态码为 200
	//c.JSON(http.StatusOK, articles)
	c.JSON(200, gin.H{"data": articles, "msg": "发送成功"})
}
