package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
)

// InsertFollow 向中间表中插入关注与被关注信息
func InsertFollow(db *gorm.DB, followerId uint, followedId uint) error {
	// 构建要插入的数据
	followData := map[string]interface{}{
		"follower_id": followerId,
		"followed_id": followedId,
	}

	// 插入数据
	result := db.Table("sw_user_follows").Create(followData)

	// 检查是否有错误
	if result.Error != nil {
		return fmt.Errorf("InsertFollow() err: %v", result.Error)
	}
	// 检查是否插入成功
	if result.RowsAffected == 0 {
		return fmt.Errorf("InsertFollow() err: 插入关注记录失败")
	}

	return nil
}

// QueryLastUserVerifyCodeByUserID 根据 UserID 查询最后一条 UserVerifyCode 记录（不论该数据的DeleteAt是否已经被赋值）
func QueryLastUserVerifyCodeByUserID(db *gorm.DB, userID uint) (*models.UserVerifyCode, error) {
	var userVerifyCode models.UserVerifyCode
	// 查询最后一条创建的记录（Unscoped()：不论这条数据的deleteAt是否被赋值）
	result := db.Unscoped().Where("user_id = ?", userID).Order("created_at DESC").First(&userVerifyCode)
	if result.Error != nil {
		return nil, fmt.Errorf("QueryLastUserVerifyCodeByUserID() err: 查找id为%d用户验证码失败, 执行的查询语句为: %v", userID, result.Statement.SQL.String())
	}

	return &userVerifyCode, nil
}

// QueryUserRank 查询用户排行。根据文章数量和文章热度计算每个用户的热度得分，并按热度从高到低排序。
// 作者榜单：1篇文章=2个热度，作者热度 = 文章数量 + 文章热度
func QueryUserRank(db *gorm.DB, page int, limit int) ([]*models.User, error) {
	var users []*models.User
	// 查询语句：根据用户的热度分数查询
	query := db.Model(&models.User{}).Order("heat DESC")

	if page > 0 && limit > 0 {
		// 计算偏移量 (从第几条记录开始查询)
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("QueryUserRank() err: %v", err)
	}

	return users, nil
}
