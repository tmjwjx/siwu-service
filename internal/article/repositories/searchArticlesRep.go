package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// QueryArticlesRep 搜索文章
func QueryArticlesRep(db *gorm.DB, req requests.ReqSearch) (articles []models.Article, err error) {
	var condition string
	if req.Kind == 0 { // 0 代表按照热度排序
		condition = "heat DESC"
	} else if req.Kind == 1 { // 1 代表按照发布时间排序
		condition = "published_at DESC"
	}

	query := db.Model(&models.Article{})

	if req.UserId != 0 { // 按用户ID筛选
		query = query.Joins("JOIN user_follows ON articles.user_id = user_follows.followed_id").
			Where("user_follows.follower_id = ?", req.UserId)
		condition = "published_at DESC" // 按照 关注 查询只能按照 时间 排序
	} else if req.Query != "" { // 按关键词搜索
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	} else if req.Tag != "" { //按标签筛选
		var tag models.Tag
		// 根据标签名称查找 tag_id
		if err = db.Where("name = ?", req.Tag).First(&tag).Error; err != nil {
			globals.Log.Errorf("err = %s", err)
			return nil, err // 返回 nil 和错误信息
		}
		fmt.Println(tag)
		query = query.Joins("JOIN article_tags ON articles.id = article_tags.article_id").
			Where("article_tags.tag_id = ?", tag.ID)
	}

	if req.CategoryId != 0 { // 按类目筛选
		query = query.Where("category_id = ?", req.CategoryId)
	}

	// 选择排序方式 时间or热度
	// 先排序 后分页
	query = query.Order(condition)

	// 判断是否分页
	if req.Limit != 0 {
		offset := (req.Page - 1) * req.Limit // 计算当前页的偏移量，用于分页
		query = query.Limit(req.Limit).Offset(offset)

	}

	// 执行查询
	if err = query.Find(&articles).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return // 结束函数执行
	}

	return articles, err
}