package repositories

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/utils"
	"gorm.io/gorm"
)

// InitUserInfoRep
// @Description: 初始化用户信息
// @Author wangyulong 2024-10-09 16:00:48
func InitUserInfoRep(db *gorm.DB, qid string, gid string) (*requests.InitUserInfoRes, error) {

	middleInfo := requests.MiddleInfo{}

	// 查询 昵称，个签
	err := db.Model(&models.User{}).
		Select("sw_users.created_at AS date, sw_users.nickname, sw_user_details.signature").
		Joins("left join sw_user_details on sw_user_details.user_id = sw_users.id").
		Where("sw_users.id = ?", gid).Scan(&middleInfo).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询 头像，昵称，个签 失败 -> %s", err)
	}

	initUserInfoRes := &requests.InitUserInfoRes{
		Nickname:  middleInfo.Nickname,
		Signature: middleInfo.Signature,
	}

	Uid, err := utils.ChangeStringToUint(gid)
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> ChangeStringToUint异常 -> %s", err)
	}

	// 查询用户头像
	paths, err := internalUtils.GetImages(db, globals.UserHome, Uid)
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询用户头像异常 -> %s", err)
	}

	for _, path := range *paths {
		initUserInfoRes.HeadShot = path
	}

	// 转化时间格式
	initUserInfoRes.Date = internalUtils.TimeFormat(&middleInfo.Date)

	// 查询
	var articleInfo []requests.ArticleInfo
	// Select("likes_count, views_count, collections_count").
	err = db.Model(&models.Article{}).Select("LikesCount", "ViewsCount", "CollectionsCount").Where("user_id = ?", gid).Scan(&articleInfo).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询 头像，昵称，个签 失败 -> %s", err)
	}

	// 统计用户总的文章数量
	initUserInfoRes.ArticleCount = len(articleInfo)

	// 统计该用户所有文章总的点赞数、浏览量、收藏数
	//var likesCount int
	//var ViewsCount int
	//var CollectionsCount int

	for _, article := range articleInfo {
		initUserInfoRes.LikesCount += article.LikesCount
		initUserInfoRes.ViewsCount += article.ViewsCount
		initUserInfoRes.CollectionsCount += article.CollectionsCount
	}

	var followed []uint // 本页面用户关注的用户
	var follower []uint // 关注本页面用户的用户
	// 查询作者关注了哪些用户
	err = db.Model(&models.UserFollow{}).Where("follower_id = ?", gid).Pluck("followed_id", &followed).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询用户关注了哪些用户失败 -> %s", err)
	}
	initUserInfoRes.ConcernsCount = len(followed)
	// 查询哪些用户关注了作者
	err = db.Model(&models.UserFollow{}).Where("followed_id = ?", gid).Pluck("follower_id", &follower).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询用户关注了哪些用户失败 -> %s", err)
	}
	initUserInfoRes.FansCount = len(follower)

	initUserInfoRes.Tag = "思悟"

	if qid == gid {
		initUserInfoRes.ConcernStatus = 2
	} else {
		// 将string类型的值转换成uint类型
		qid2, err := utils.ChangeStringToUint(qid)
		if err != nil {
			return nil, fmt.Errorf("InitUserInfoRep -> %s", err)
		}

		for _, id := range follower {
			if id == qid2 {
				initUserInfoRes.ConcernStatus = 1
				break
			}
		}
		//initUserInfoRes.ConcernStatus = 0
	}

	return initUserInfoRes, nil
}

// EditSignatureRep
// @Description: 编辑个签
// @Author wangyulong 2024-10-11 16:09:26
// @param        db *gorm.DB
// @param        req *requests.EditSignatureReq
// @param        id uint
// @return       error
func EditSignatureRep(db *gorm.DB, req *requests.EditSignatureReq, id uint) error {

	userDetail := &models.UserDetail{
		UserID: id,
	}

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("EditSignatureRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询该用户详细信息是否存在
	err := tx.First(userDetail).Error
	if err != nil {
		// 该用户详细信息不存在，插入个签
		userDetail.Signature = req.Signature
		err = tx.Model(&models.UserDetail{}).Create(userDetail).Error
		if err != nil {
			// 回滚事务
			tx.Rollback()
			return fmt.Errorf("EditSignatureRep -> 插入个签失败 -> %s", err)
		}
	} else {
		// 该用户详细信息存在，只更新 个签 这一个字段
		err = tx.Model(&models.UserDetail{}).Where("user_id = ?", id).Update("signature", req.Signature).Error
		if err != nil {
			// 回滚事务
			tx.Rollback()
			return fmt.Errorf("EditSignatureRep -> 更新个签失败 -> %s", err)
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("EditSignatureRep -> 提交事务失败 -> %s", err)
	}

	return nil
}
