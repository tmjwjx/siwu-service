package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

func SearchArticles(db *gorm.DB, req requests.SearchRequest, offset int) (articles []models.Article, err error) {
	var condition string
	if req.Kind == 0 { // 0 代表按照热度排序
		condition = "heat DESC"
	} else if req.Kind == 1 { // 1 代表按照发布时间排序
		condition = "published_at DESC"
	}
	
	query := db.Model(&models.Article{})
	
	// 按用户ID筛选
	if req.UserId != 0 {
		query = query.Joins("JOIN user_follows ON articles.user_id = user_follows.followed_id").
			Where("user_follows.follower_id = ?", req.UserId)
		//condition = "published_at DESC" // 按照时间排序
	}
	// 按关键词搜索
	if req.Query != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	}
	// 按类目筛选
	if req.CategoryId != 0 {
		query = query.Where("category_id = ?", req.CategoryId)
	}
	//// 按标签筛选
	fmt.Println(req)
	fmt.Println(req.Tag)
	if req.Tag != "" {
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
	
	query = query.Order(condition).Limit(req.Limit).Offset(offset)
	
	// 执行查询
	if err = query.Find(&articles).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return // 结束函数执行
	}
	
	return articles, err
}