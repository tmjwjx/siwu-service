package repositories

import (
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
	"time"
)

// InsertArticlesRep 新建文章
func InsertArticlesRep(db *gorm.DB, req requests.ReqPublish) error {

	newArticle := models.Article{
		UserID:           req.UserId,
		Title:            req.Title,
		LikesCount:       0,
		CollectionsCount: 0,
		CommentsCount:    0,
		ViewsCount:       0,
		Heat:             0,
		Status:           req.Status,
		CategoryID:       req.CategoryID,
		Summary:          req.Summary,
		Content:          req.Content,
	}
	result := db.Create(&newArticle)
	if result.Error != nil {
		globals.Log.Fatal("创建文章失败:", result.Error)
	}

	return nil
}

// UpdatePublishRep 设置文章的发布时间
func UpdatePublishRep(db *gorm.DB, article *models.Article) error {

	article.PublishedAt = time.Now()

	// 更新数据库中的文章记录
	if err := db.Save(article).Error; err != nil {
		return err // 返回错误
	}

	return nil
}

// QueryTag 返回全部标签
func QueryTag(db *gorm.DB) (tags []models.Tag, err error) {
	// 查询 Tag 表中的所有数据
	if err = db.Find(&tags).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	return tags, nil
}

// QueryCategory 返回全部类目
func QueryCategory(db *gorm.DB) (categories []models.Category, err error) {

	// 查询 Category 表中的所有数据
	if err = db.Find(&categories).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	return categories, nil

}