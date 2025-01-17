package repositories

import (
	"fmt"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"gorm.io/gorm"
)

// QueryLastUserVerifyCodeByEmail 根据 email 查询最后一条 UserVerifyCode 记录（不论该数据的DeleteAt是否已经被赋值）
func QueryLastUserVerifyCodeByEmail(db *gorm.DB, email string) (*models.UserVerifyCode, error) {
	var userVerifyCode models.UserVerifyCode
	// 查询最后一条创建的记录（Unscoped()：不论这条数据的deleteAt是否被赋值）
	result := db.Unscoped().Where("email = ?", email).Order("created_at DESC").First(&userVerifyCode)
	if result.Error != nil {
		return nil, fmt.Errorf("QueryLastUserVerifyCodeByUserID() err: 查找email为%s用户验证码失败, 执行的查询语句为: %v", email, result.Statement.SQL.String())
	}
	return &userVerifyCode, nil
}

// QueryUserRank 查询用户排行。根据文章数量和文章热度计算每个用户的热度得分，并按热度从高到低排序。
// 作者榜单：1篇文章=2个热度，作者热度 = 文章数量 + 文章热度
func QueryUserRank(db *gorm.DB, page int, limit int) ([]*models.User, error) {
	var users []*models.User
	// 查询语句：根据用户的热度分数查询。Paginate：可复用的查询逻辑，对查询结果进行分页。
	query := db.Model(&models.User{}).Order("heat DESC").Scopes(sqlUtils.Paginate(page, limit))

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("QueryUserRank() err: %v", err)
	}

	return users, nil
}

// QueryAttentionByPage
// @Description: 搜索用户关注的人。
// @Author lizhuang 2025-01-17 14:15:42
// @param        db *gorm.DB
// @param        keyword string 模糊查询昵称
// @param        page
// @param        limit int
func QueryAttentionByPage(db *gorm.DB, followerId uint, keyword string, page, limit int) ([]uint, error) {
	var ids []uint

	// 构造查询
	query := db.Model(&models.UserFollow{}).
		Scopes(sqlUtils.Paginate(page, limit)). // 分页
		Where("follower_id = ?", followerId).   // 根据关注者id查询
		Pluck("followed_id", &ids)              // 查询被关注者

	// 执行查询
	if err := query.Error; err != nil {
		return nil, fmt.Errorf("QueryAttentionByPage() err: %v", err)
	}

	// 模糊查询来筛选 nickname
	newIds, err := QueryUserIdsByNickname(db, ids, keyword)
	if err != nil {
		return nil, fmt.Errorf("QueryAttentionByPage() -> %v", err)
	}

	return newIds, nil
}

// QueryUserIdsByNickname 使用模糊查询来查询 nickname 字段里面包含 keyword 的内容。
func QueryUserIdsByNickname(db *gorm.DB, ids []uint, keyword string) ([]uint, error) {
	var userIds []uint

	// 执行查询：根据id在给定的ids切片中，并且nickname字段包含keyword
	err := db.Model(&models.User{}).
		Where("id IN (?)", ids).                   // 查询指定id的用户
		Where("nickname LIKE ?", "%"+keyword+"%"). // 使用LIKE进行模糊查询
		Pluck("id", &userIds).Error                // 只查询并返回用户id字段

	if err != nil {
		return nil, fmt.Errorf("QueryUserIdsByNickname() err: %v", err)
	}
	return userIds, nil
}
