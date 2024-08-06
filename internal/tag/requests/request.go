package requests

import "gorm.io/gorm"

// Tag 标签
type Tag struct {
	gorm.Model          //ID CreatedAt UpdatedAt DeletedAt
	Name         string `form:"name"`          // 标签名称
	Description  string `form:"description"`   // 标签描述
	ArticleCount int    `form:"article_count"` // 标签关联的文章数量
	Heat         int    `form:"heat"`          // 标签热度
	FansCount    uint   `form:"fan_count"`     // 关注人数
	Path         string `form:"path"`          // 相关图片
}
