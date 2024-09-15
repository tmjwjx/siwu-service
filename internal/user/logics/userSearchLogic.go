package logics

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"gorm.io/gorm"
)

// Follow 关注和取消关注。followerId 关注 followedId
func (u *UserReqContext) Follow(follow requests.FollowMsg) error {
	followerId := follow.FollowerId
	followedId := follow.FollowedId

	// 查询两个id，判断两个id是否存在
	// 关注者
	user1 := repositories.QueryUserById(u.DB, followerId)
	if user1 == nil {
		return fmt.Errorf("UserReqContext.Follow() : 关注者不存在")
	}
	// 被关注者
	user2 := repositories.QueryUserById(u.DB, followedId)
	if user2 == nil {
		return fmt.Errorf("UserReqContext.Follow() : 被关注者不存在")
	}

	// 判断是否已经关注过了，如果已经关注过了，再次点击就会取消关注
	// 我的关注
	followedIDSli, err := repositories.QueryFollowed(u.DB, followerId)
	if err != nil {
		return fmt.Errorf("UserReqContext.Follow() -> %v: ", err)
	}
	// 是否已经关注过 followedId
	var isFollowed = false
	for _, id := range followedIDSli {
		if id == followedId {
			isFollowed = true
			break
		}
	}

	// followedId的粉丝
	followerIDSli, err := repositories.QueryFollower(u.DB, followedId)
	if err != nil {
		return fmt.Errorf("UserReqContext.Follow() -> %v: ", err)
	}

	// 关注
	if !isFollowed {
		if err := repositories.InsertFollow(u.DB, followerId, followedId); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		//  followerId关注数量+1，followedId粉丝数量+1
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) + 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) + 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}

	} else { // 取消关注
		if _, err := repositories.DeleteObjectsByTable(u.DB, "sw_user_follows", map[string]interface{}{"follower_id": followerId, "followed_id": followedId}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		//  followerId关注数量-1，followedId粉丝数量-1
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) - 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) - 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
	}

	return nil
}

// UserRank 用户热度排行
func (u *UserReqContext) UserRank(msg requests.UserRankMsg) ([]*models.User, error) {
	// 查询
	userRank, err := repositories.QueryUserRank(u.DB, msg.Page, msg.Limit)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
	}

	// 选择需要的数据

	return userRank, nil
}
