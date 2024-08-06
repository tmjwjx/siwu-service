package repositories

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
)

// UpdateTagUserCountReq 更新数据库中标签的关注人数
func UpdateTagUserCountReq(tagID uint) (string, error) {
	var tag models.Tag
	var fansCount string // 统计现在的人数
	if err := globals.DB.Take(&tag, "id = ?", tagID).Error; err != nil {
		return "", err
	}
	if err := globals.DB.Model(&tag).Update("fans_count", tag.FansCount+1).Error; err != nil {
		return "", err
	}
	if tag.FansCount > 1000 {
		fansCount = fmt.Sprintf("%.1fk", tag.FansCount/1000)
	} else {
		fansCount = fmt.Sprintf("%d", tag.FansCount)
	}
	return fansCount, nil
}

// UpdateTagArticleCountReq 更新数据库中标签的文章数量
func UpdateTagArticleCountReq() {
	//var article []models.Article

}
