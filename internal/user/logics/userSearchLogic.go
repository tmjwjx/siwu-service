package logics

import (
	"fmt"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// Follow 关注和取消关注。followerId 关注 followedId
func (u *UserReqContext) Follow(follow requests.FollowReq) error {
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
		userFollow := &models.UserFollow{
			FollowerId: followerId,
			FollowedId: followedId,
		}
		// 插入关注数据
		if err := sqlUtils.InsertObject(u.DB, userFollow); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}

		//  followerId关注数量+1，followedId粉丝数量+1
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) + 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) + 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}

	} else { // 取消关注
		if _, err := sqlUtils.DeleteObjectsByTable(u.DB, "sw_user_follows", map[string]interface{}{"follower_id": followerId, "followed_id": followedId}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		//  followerId关注数量-1，followedId粉丝数量-1
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) - 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) - 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
	}

	return nil
}

// UserRank 用户热度排行
func (u *UserReqContext) UserRank(id uint, msg requests.UserRankReq) ([]*requests.UserRankRes, error) {
	userRankReqSli := make([]*requests.UserRankRes, msg.Limit)

	// 查询排行榜：每个用户id，昵称
	usersRank, err := repositories.QueryUserRank(u.DB, msg.Page, msg.Limit)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
	}

	// 查询每个用户是否已关注（用户是否关注了这个排行榜上的用户：用户是否关注了这个排行榜上的用户：未关注：0，已关注：1，这个用户是自己：2）
	// 查询 id 关注了谁
	ids, err := repositories.QueryFollowed(u.DB, id)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
	}

	// 使用 lo.Associate 将 []uint 转换为 map[uint]bool
	m := lo.Associate(ids, func(item uint) (uint, bool) {
		return item, true
	})

	// 查询每个id的详细信息
	for i, v := range usersRank {
		// 查询用户的职业描述
		userDetail := repositories.QueryUserDetailsById(u.DB, v.ID)
		if userDetail == nil {
			return nil, fmt.Errorf("UserReqContext.UserRank() err: user_id为 %v 的userDetail不存在", v.ID)
		}
		userRankReqSli[i] = &requests.UserRankRes{
			Id:              v.ID,
			Nickname:        v.Nickname,
			CareerDirection: userDetail.CareerDirection,
		}

		// 查询用户的头像路径
		//userImgs, err := controllers.GetImagesControllers("用户", v.ID)
		//if err != nil {
		//	// 数据库中没有该用户的头像，使用默认的头像
		//	userRankReqSli[i].AvatarPath = internalUtils.UserDefaultImage
		//} else {
		//	userRankReqSli[i].AvatarPath = (*userImgs)[0].Path
		//}

		// 未关注：0，已关注：1，这个用户是自己：2
		if v.ID == id {
			userRankReqSli[i].IsFollowed = 2
		} else {
			// 判断用户是否关注该id
			_, ok := m[v.ID]
			if !ok {
				userRankReqSli[i].IsFollowed = 0
			} else {
				userRankReqSli[i].IsFollowed = 1
			}
		}
	}

	return userRankReqSli, nil
}
