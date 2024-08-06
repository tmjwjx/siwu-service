package repositorys

import (
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchArticles(c *gin.Context, req requests.SearchRequest) (articles []models.Article) {
	db := globals.DB
	offset := (req.Page - 1) * req.Limit // 计算当前页的偏移量，用于分页
	var condition string
	if req.Kind == 0 { // 0 代表按照热度排序
		condition = "heat DESC"
	} else if req.Kind == 1 { // 1 代表按照发布时间排序
		condition = "published_at DESC"
	}

	if err := db.Where("(title LIKE ? OR summary LIKE ?) AND category_id = ?", "%"+req.Query+"%", "%"+req.Query+"%", req.Category).
		Order(condition).                   // 按照热度降序排序
		Limit(req.Limit).                   // 限制返回的产品数量
		Offset(offset).                     // 设置查询的偏移量
		Find(&articles).Error; err != nil { // 执行查询并检查是否出错
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Database errors"}) // 返回 500 错误
		return                                                                     // 结束函数执行
	}

	return articles
}
