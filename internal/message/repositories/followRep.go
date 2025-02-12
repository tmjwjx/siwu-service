package repositories

import (
	"errors"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// FollowUnreadCount
// @Description: 关注消息未读数量
// @param        db *gorm.DB
// @param        id uint
// @return       count
// @return       err
// @Author tianjiajie 2025-02-12 21:25:00
func FollowUnreadCount(db *gorm.DB, id uint) (count int64, err error) {
	err = db.Model(&models.UserFollow{}).
		Where("followed_id = ?", id).
		Where("is_read = ?", 0).
		Count(&count).
		Error

	if err != nil {
		globals.Log.Errorf("Failed to count unread likes: %v", err)
		return 0, err
	}

	return count, nil
}

// FollowRead
// @Description: 关注消息已读
// @param        db *gorm.DB
// @param        id uint
// @return       err
// @Author tianjiajie 2025-02-12 21:11:52
func FollowRead(db *gorm.DB, id uint) (err error) {
	err = db.Model(&models.UserFollow{}).
		Where("followed_id = ?", id).
		Update("is_read", 1).
		Error

	if err != nil {
		globals.Log.Errorf("Failed to mark likes as read: %v", err)
		return err
	}

	return nil
}

// FollowRep
// @Description: 关注消息
// @param        *gorm.DB *gorm.DB
// @param        requests.MessageReq requests.MessageReq
// @param        uint uint
// @return       res
// @return       err
// @Author tianjiajie 2025-01-16 20:33:05
func FollowRep(db *gorm.DB, req requests.MessageReq, userId uint) (res []requests.FollowMessageRes, err error) {

	// 获取关注消息
	err = db.Model(&models.UserFollow{}).
		Joins("join sw_users on sw_users.id = sw_user_follows.follower_id").
		Select("sw_user_follows.follower_id, sw_users.nickname, sw_user_follows.created_at").
		Where("sw_user_follows.followed_id = ?", userId).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("sw_user_follows.created_at DESC").
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	// 遍历res，查询用户的头像路径 和 是否关注
	for i := 0; i < len(res); i++ {
		// 查询用户的头像路径
		userImages, err := internalUtils.GetImages(db, globals.UserHome, res[i].FollowerId)
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.GetInfo() %v", err)
		}
		// 没有图片
		if userImages == nil {
			return nil, fmt.Errorf("UserReqContext.GetInfo() err = 无法找到id为%d的用户图片", res[i].FollowerId)
		}
		res[i].Path = (*userImages)[0]

		// 查询是否关注
		var follow int64
		err = db.Model(&models.UserFollow{}).
			Where("follower_id = ? and followed_id = ?", userId, res[i].FollowerId).
			Count(&follow).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("follow", follow)
		if follow != 0 {
			res[i].IsFollowed = 1
		} else {
			res[i].IsFollowed = 0
		}
	}

	// 格式化时间
	for i := range res {
		res[i].FormatTime = internalUtils.TimeFormat(res[i].CreatedAt)
		res[i].DailyTime = internalUtils.TimeFormatDaily(res[i].CreatedAt)
	}

	// 关注消息已读
	err = FollowRead(db, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
