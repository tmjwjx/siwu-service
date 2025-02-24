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
// @Description: 分页搜索用户关注的人。
// @Author lizhuang 2025-01-17 14:15:42
// @param        db *gorm.DB
// @param        keyword string 模糊查询昵称
// @param        page
// @param        limit int
// func QueryAttentionByPage(db *gorm.DB, followerId uint, keyword string, page, limit int) ([]uint, error) {
// 	var ids []uint
//
// 	// 构造查询
// 	// query := db.Model(&models.UserFollow{}).
// 	// 	Scopes(sqlUtils.Paginate(page, limit)). // 分页
// 	// 	Where("follower_id = ?", followerId).   // 根据关注者id查询
// 	// 	Pluck("followed_id", &ids)              // 查询被关注者
// 	//
// 	// // 执行查询
// 	// if err := query.Error; err != nil {
// 	// 	return nil, fmt.Errorf("QueryAttentionByPage() err: %v", err)
// 	// }
// 	//
// 	// // 模糊查询来筛选 nickname
// 	// newIds, err := QueryUserIdsByNickname(db, ids, keyword)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("QueryAttentionByPage() -> %v", err)
// 	// }
// 	//
// 	// return newIds, nil
//
// 	query := db.Model(&models.UserFollow{}).
// 		Scopes(sqlUtils.Paginate(page, limit)). // 分页
// 		Where("follower_id = ?", followerId)    // 根据关注者id查询
//
// 	if keyword != "" {
// 		// 先通过 UserFollow 表关联到 User 表，再进行模糊查询
// 		query = query.Joins("JOIN sw_users ON sw_user_follows.followed_id = sw_users.id").
// 			Where("sw_users.nickname LIKE ?", "%"+keyword+"%")
// 	}
//
// 	// 查询被关注者的ID
// 	if err := query.Pluck("sw_user_follows.followed_id", &ids).Error; err != nil {
// 		return nil, fmt.Errorf("QueryAttentionByPage() err: %v", err)
// 	}
//
// 	return ids, nil
// }

// QueryAttentionByPage
// @Description: 分页搜索用户关注的人。
// @Author lizhuang 2025-01-17 14:15:42
// @param        db *gorm.DB
// @param        followerId uint 关注者ID
// @param        keyword string 模糊查询昵称
// @param        page int 第几页
// @param        limit int 每页条数
// @return       []uint 被关注者的ID列表
// @return       bool 是否还有更多数据
// @return       error 错误信息
func QueryAttentionByPage(db *gorm.DB, followerId uint, keyword string, page, limit int) ([]uint, bool, error) {
	var ids []uint

	query := db.Model(&models.UserFollow{}).
		Where("follower_id = ?", followerId) // 根据关注者id查询

	if keyword != "" {
		// 关联 User 表进行昵称模糊查询
		query = query.Joins("JOIN sw_users ON sw_user_follows.followed_id = sw_users.id").
			Where("sw_users.nickname LIKE ?", "%"+keyword+"%")
	}

	// 先获取总数，用于判断是否还有更多数据
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, false, fmt.Errorf("QueryAttentionByPage() count err: %v", err)
	}

	// 添加分页并查询结果
	if err := query.Scopes(sqlUtils.Paginate(page, limit)).
		Pluck("sw_user_follows.followed_id", &ids).Error; err != nil {
		return nil, false, fmt.Errorf("QueryAttentionByPage() err: %v", err)
	}

	// 判断是否还有更多数据
	// 当前偏移量 = (page-1) * limit
	// 如果偏移量 + 当前页数据量 < 总数，则还有数据
	hasMore := (page-1)*limit+len(ids) < int(total)

	return ids, hasMore, nil
}

// QueryAttention
// @Description: 搜索用户关注的全部人
// @Author lizhuang 2025-01-22 09:08:25
// @param        db *gorm.DB
// @param        followerId uint
// @return       []uint
// @return       error
func QueryAttention(db *gorm.DB, followerId uint) ([]uint, error) {
	var ids []uint

	// 构造查询
	query := db.Model(&models.UserFollow{}).
		Where("follower_id = ?", followerId). // 根据关注者id查询
		Pluck("followed_id", &ids)            // 查询被关注者

	// 执行查询
	if err := query.Error; err != nil {
		return nil, fmt.Errorf("QueryAttentionByPage() err: %v", err)
	}

	return ids, nil
}

// QueryUserIDArticleNumOfPub
// @Description: 根据要查询的文章status，查询用户的文章数量
// @Author lizhuang 2025-01-22 09:17:39
// @param        db *gorm.DB
// @param        userId uint
// @param        status string
// @return       int64
// @return       error
func QueryUserIDArticleNumOfPub(db *gorm.DB, userId uint, status string) (int64, error) {
	// 构造查询
	query := db.Model(&models.Article{}).
		Where("user_id = ?", userId). // 根据关注者id查询
		Where("status = ?", status)   // 要查询的文章状态

	var count int64
	// 执行查询
	if err := query.Count(&count).Error; err != nil {
		return -1, fmt.Errorf("QueryUserIDArticleNumOfPub() err: %v", err)
	}

	return count, nil
}

// // QueryUserIdsByNickname 使用模糊查询来查询 nickname 字段里面包含 keyword 的内容。
// func QueryUserIdsByNickname(db *gorm.DB, ids []uint, keyword string) ([]uint, error) {
// 	var userIds []uint
//
// 	// 执行查询：根据id在给定的ids切片中，并且nickname字段包含keyword
// 	err := db.Model(&models.User{}).
// 		Where("id IN (?)", ids).                   // 查询指定id的用户
// 		Where("nickname LIKE ?", "%"+keyword+"%"). // 使用LIKE进行模糊查询
// 		Pluck("id", &userIds).Error                // 只查询并返回用户id字段
//
// 	if err != nil {
// 		return nil, fmt.Errorf("QueryUserIdsByNickname() err: %v", err)
// 	}
// 	return userIds, nil
// }

// QueryUserArticleRep
// @Description: 查询用户的文章
// @param        db *gorm.DB
// @param        req requests.UserDataRequest
// @param        b bool
// @return       []*models.Article
// @return       error
// @Author tianjiajie 2025-01-17 17:19:16
// func QueryUserArticleRep(db *gorm.DB, req requests.UserDataRequest, b bool) ([]article_req.SearchArticleListRes, int64, error) {
//	var articles []article_req.SearchArticleListRes
//	query := db.
//		Model(&models.Article{}).Preload("Tags").
//		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id").
//		//Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id").
//		Select("DISTINCT sw_articles.*,sw_articles.published_at, sw_users.nickname")
//
//	// 查询条件
//	query = query.Where("user_id = ?", req.Id)
//
//	// 排序
//	query = query.Order("sw_articles.published_at DESC")
//
//	// 如果不是自己，只查询已发布的文章
//	if !b {
//		query = query.Where("sw_articles.status = ?", "public")
//	}
//
//	// 查询总数
//	var total int64
//	// 复制查询语句
//	if err := query.Select("sw_articles.id").Count(&total).Error; err != nil {
//		return nil, 0, fmt.Errorf("QueryUserArticleRep() -> %v", err)
//	}
//
//	// 分页
//	query = query.Scopes(sqlUtils.Paginate(req.Page, req.Limit))
//
//	// 查询文章
//	if err := query.Find(&articles).Error; err != nil {
//		return nil, 0, fmt.Errorf("QueryUserArticleRep() -> %v", err)
//	}
//	return articles, total, nil
// }
