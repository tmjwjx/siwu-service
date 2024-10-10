package repositories

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/internal/user/requests"
	"gorm.io/gorm"
)

// InitUserInfoRep
// @Description: 初始化用户信息
// @Author wangyulong 2024-10-09 16:00:48
func InitUserInfoRep(db *gorm.DB, qid string, gid string) (*requests.InitUserInfoRes, error) {

	initUserInfoRes := &requests.InitUserInfoRes{}
	// 查询 头像，昵称，个签
	err := db.Model(&models.User{}).Select("sw_users.created_at AS data, sw_users.nickname, sw_user_details.signature, sw_attachments.path").Joins("join sw_user_details on sw_user_details.user_id = sw_users.id").
		Joins("join sw_attachments on sw_attachments.home_id = sw_users.id").Where("sw_users.id = ?", gid).Scan(initUserInfoRes).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询 头像，昵称，个签 失败 -> %s", err)
	}
	// 查询
	var articleInfo []requests.ArticleInfo
	// Select("likes_count, views_count, collections_count").
	err = db.Model(&models.Article{}).Select("LikesCount", "ViewsCount", "CollectionsCount").Where("user_id = ?", gid).Scan(&articleInfo).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询 头像，昵称，个签 失败 -> %s", err)
	}
	// 统计该用户所有文章总的点赞数、浏览量、收藏数
	//var likesCount int
	//var ViewsCount int
	//var CollectionsCount int

	for _, article := range articleInfo {
		initUserInfoRes.LikesCount += article.LikesCount
		initUserInfoRes.ViewsCount += article.ViewsCount
		initUserInfoRes.CollectionsCount += article.CollectionsCount
	}

	var followed []int
	var follower []int
	// 查询作者关注了哪些用户
	err = db.Model(&models.UserFollow{}).Select("followed_id").Where("follower_id = ?", gid).Find(&followed).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询用户关注了哪些用户失败 -> %s", err)
	}
	initUserInfoRes.ConcernsCount = len(followed)
	// 查询哪些用户关注了作者
	err = db.Model(&models.UserFollow{}).Select("follower_id").Where("followed_id = ?", gid).Find(&follower).Error
	if err != nil {
		return nil, fmt.Errorf("InitUserInfoRep -> 查询用户关注了哪些用户失败 -> %s", err)
	}
	initUserInfoRes.FansCount = len(follower)

	initUserInfoRes.Tag = "思悟"

	if qid == gid {
		initUserInfoRes.ConcernStatus = 2
	} else {
		// 将string类型的值转换成int类型
		gid2, err := internalUtils.ChangeStringToInt(gid)
		if err != nil {
			return nil, fmt.Errorf("InitUserInfoRep -> %s", err)
		}

		for _, id := range followed {
			if id == gid2 {
				initUserInfoRes.ConcernStatus = 1
				break
			}
		}
		initUserInfoRes.ConcernStatus = 0
	}

	return initUserInfoRes, nil
}
