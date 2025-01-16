package repositories

import (
	"fmt"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"gorm.io/gorm"
)

// // QueryLastUserVerifyCodeByUserID 根据 UserID 查询最后一条 UserVerifyCode 记录（不论该数据的DeleteAt是否已经被赋值）
// func QueryLastUserVerifyCodeByUserID(db *gorm.DB, userID uint) (*models.UserVerifyCode, error) {
// 	var userVerifyCode models.UserVerifyCode
// 	// 查询最后一条创建的记录（Unscoped()：不论这条数据的deleteAt是否被赋值）
// 	result := db.Unscoped().Where("user_id = ?", userID).Order("created_at DESC").First(&userVerifyCode)
// 	if result.Error != nil {
// 		return nil, fmt.Errorf("QueryLastUserVerifyCodeByUserID() err: 查找id为%d用户验证码失败, 执行的查询语句为: %v", userID, result.Statement.SQL.String())
// 	}
//
// 	return &userVerifyCode, nil
// }

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
