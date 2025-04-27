package logics

import (
	"encoding/json"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"strconv"
)

// ClickAttention 点击关注和点击取消关注。followerId 关注 followedId
func (u *UserReqContext) ClickAttention(req requests.ClickAttentionReq) (int, error) {
	followerId := req.FollowerId
	followedId := req.FollowedId

	// 查询两个id，判断两个id是否存在
	// 关注者
	user1 := repositories.QueryUserById(u.DB, followerId)
	if user1 == nil {
		// return fmt.Errorf("UserReqContext.ClickAttention() : 关注者不存在")
		globals.Log.Error(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(followerId)))
		return 400, fmt.Errorf(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(followerId)))
	}
	// 被关注者
	user2 := repositories.QueryUserById(u.DB, followedId)
	if user2 == nil {
		// return fmt.Errorf("UserReqContext.ClickAttention() : 被关注者不存在")
		globals.Log.Error(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(followedId)))
		return 400, fmt.Errorf(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(followedId)))
	}

	// 判断是否已经关注过了，如果已经关注过了，再次点击就会取消关注
	// 我的关注
	followedIDSli, err := repositories.QueryFollowed(u.DB, followerId)
	if err != nil {
		// return fmt.Errorf("UserReqContext.ClickAttention() -> %v: ", err)
		globals.Log.Error(err.Error())
		return 500, err
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
		// return fmt.Errorf("UserReqContext.ClickAttention() -> %v: ", err)
		globals.Log.Error(err.Error())
		return 500, err
	}

	// 关注
	if !isFollowed {
		userFollow := &models.UserFollow{
			FollowerId: followerId,
			FollowedId: followedId,
		}
		// 插入关注数据
		if err = sqlUtils.InsertObject(u.DB, userFollow); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}

		//  followerId关注数量+1，followedId粉丝数量+1
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) + 1}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) + 1}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}

		// user2 的热度+10
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"heat": user2.Heat + 10}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}

	} else { // 取消关注
		if _, err := sqlUtils.DeleteObjectsByTable(u.DB, "sw_user_follows", map[string]interface{}{"follower_id": followerId, "followed_id": followedId}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}
		//  followerId关注数量-1，followedId粉丝数量-1
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) - 1}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) - 1}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}

		// user2 的热度-10
		if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"heat": user2.Heat - 10}); err != nil {
			// return fmt.Errorf("UserReqContext.ClickAttention() -> %v", err)
			globals.Log.Error(err.Error())
			return 500, err
		}
	}

	// 通知用户关注消息
	//internalUtils.MessagePush("follow", strconv.Itoa(int(followedId)))
	jsonData, _ := json.Marshal(gin.H{"content": "follow"})
	internalUtils.MessagePush2(string(jsonData), strconv.Itoa(int(followedId)), globals.NoticeType, "")

	return 200, nil
}

// UserRank 用户热度排行
func (u *UserReqContext) UserRank(req requests.UserRankReq) ([]*requests.UserRankRes, error) {
	// userRankReqSli := make([]*requests.UserRankRes, req.Limit)
	//
	// // 查询排行榜：每个用户id，昵称
	// usersRank, err := repositories.QueryUserRank(u.DB, req.Page, req.Limit)
	// if err != nil {
	// 	// return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
	// 	globals.Log.Error(err.Error())
	// 	return nil, err
	// }
	//
	// // 查询每个用户是否已关注（用户是否关注了这个排行榜上的用户：用户是否关注了这个排行榜上的用户：未关注：0，已关注：1，这个用户是自己：2）
	// // 查询 id 关注了谁
	// ids, err := repositories.QueryFollowed(u.DB, id)
	// if err != nil {
	// 	// return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
	// 	globals.Log.Error(err.Error())
	// 	return nil, err
	// }
	//
	// // 使用 lo.Associate 将 []uint 转换为 map[uint]bool
	// m := lo.Associate(ids, func(item uint) (uint, bool) {
	// 	return item, true
	// })
	//
	// // 查询每个id的详细信息
	// for i, v := range usersRank {
	// 	// 查询用户的职业描述
	// 	userDetail := repositories.QueryUserDetailsById(u.DB, v.ID)
	// 	if userDetail == nil {
	// 		// return nil, fmt.Errorf("UserReqContext.UserRank() err: user_id为 %v 的userDetail不存在", v.ID)
	// 		globals.Log.Error(response.ErrUserDetailNotExist + ":" + strconv.Itoa(int(v.ID)))
	// 		return nil, fmt.Errorf(response.ErrUserDetailNotExist + ":" + strconv.Itoa(int(v.ID)))
	// 	}
	// 	userRankReqSli[i] = &requests.UserRankRes{
	// 		Id:              v.ID,
	// 		Nickname:        v.Nickname,
	// 		CareerDirection: userDetail.CareerDirection,
	// 	}
	//
	// 	// 查询用户的头像路径
	// 	userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, v.ID)
	// 	if err != nil {
	// 		// return nil, fmt.Errorf("UserReqContext.UserRank() %v", err)
	// 		globals.Log.Error(err.Error())
	// 		return nil, err
	// 	}
	// 	// 没有图片
	// 	if userImages == nil {
	// 		// return nil, fmt.Errorf("UserReqContext.UserRank() err: 无法找到id为%d的用户头像图片", v.ID)
	// 		globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(v.ID)))
	// 		return nil, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(v.ID)))
	// 	}
	// 	userRankReqSli[i].AvatarPath = (*userImages)[0]
	//
	// 	// 未关注：0，已关注：1，这个用户是自己：2
	// 	if v.ID == id {
	// 		userRankReqSli[i].IsFollowed = 2
	// 	} else {
	// 		// 判断用户是否关注该id
	// 		_, ok := m[v.ID]
	// 		if !ok {
	// 			userRankReqSli[i].IsFollowed = 0
	// 		} else {
	// 			userRankReqSli[i].IsFollowed = 1
	// 		}
	// 	}
	// }
	//
	// return userRankReqSli, nil

	userRankReqSli := make([]*requests.UserRankRes, req.Limit)

	// 判断是否是游客模式
	if req.Id == 0 { // 如果 id == 0 那么就是游客模式
		// 查询排行榜：每个用户id，昵称
		usersRank, err := repositories.QueryUserRank(u.DB, req.Page, req.Limit)
		if err != nil {
			// return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
			globals.Log.Error(err.Error())
			return nil, err
		}

		// 查询每个id的详细信息
		for i, v := range usersRank {
			// 查询用户的职业描述
			userDetail := repositories.QueryUserDetailsById(u.DB, v.ID)
			if userDetail == nil {
				// return nil, fmt.Errorf("UserReqContext.UserRank() err: user_id为 %v 的userDetail不存在", v.ID)
				globals.Log.Error(response.ErrUserDetailNotExist + ":" + strconv.Itoa(int(v.ID)))
				return nil, fmt.Errorf(response.ErrUserDetailNotExist + ":" + strconv.Itoa(int(v.ID)))
			}
			userRankReqSli[i] = &requests.UserRankRes{
				Id:              v.ID,
				Nickname:        v.Nickname,
				CareerDirection: userDetail.CareerDirection,
			}

			// 查询用户的头像路径
			userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, v.ID)
			if err != nil {
				// return nil, fmt.Errorf("UserReqContext.UserRank() %v", err)
				globals.Log.Error(err.Error())
				return nil, err
			}
			// 没有图片
			if userImages == nil {
				// return nil, fmt.Errorf("UserReqContext.UserRank() err: 无法找到id为%d的用户头像图片", v.ID)
				globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(v.ID)))
				return nil, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(v.ID)))
			}
			userRankReqSli[i].AvatarPath = (*userImages)[0]

		}
		return userRankReqSli, nil

	} else {
		// 查询排行榜：每个用户id，昵称
		usersRank, err := repositories.QueryUserRank(u.DB, req.Page, req.Limit)
		if err != nil {
			// return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
			globals.Log.Error(err.Error())
			return nil, err
		}

		// 查询每个用户是否已关注（用户是否关注了这个排行榜上的用户：用户是否关注了这个排行榜上的用户：未关注：0，已关注：1，这个用户是自己：2）
		// 查询 id 关注了谁
		ids, err := repositories.QueryFollowed(u.DB, req.Id)
		if err != nil {
			// return nil, fmt.Errorf("UserReqContext.UserRank() -> %v", err)
			globals.Log.Error(err.Error())
			return nil, err
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
				// return nil, fmt.Errorf("UserReqContext.UserRank() err: user_id为 %v 的userDetail不存在", v.ID)
				globals.Log.Error(response.ErrUserDetailNotExist + ":" + strconv.Itoa(int(v.ID)))
				return nil, fmt.Errorf(response.ErrUserDetailNotExist + ":" + strconv.Itoa(int(v.ID)))
			}
			userRankReqSli[i] = &requests.UserRankRes{
				Id:              v.ID,
				Nickname:        v.Nickname,
				CareerDirection: userDetail.CareerDirection,
			}

			// 查询用户的头像路径
			userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, v.ID)
			if err != nil {
				// return nil, fmt.Errorf("UserReqContext.UserRank() %v", err)
				globals.Log.Error(err.Error())
				return nil, err
			}
			// 没有图片
			if userImages == nil {
				// return nil, fmt.Errorf("UserReqContext.UserRank() err: 无法找到id为%d的用户头像图片", v.ID)
				globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(v.ID)))
				return nil, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(v.ID)))
			}
			userRankReqSli[i].AvatarPath = (*userImages)[0]

			// 未关注：0，已关注：1，这个用户是自己：2
			if v.ID == req.Id {
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
}

// Attention 搜索用户关注的人
func (u *UserReqContext) Attention(req requests.AttentionReq) (*requests.AttentionRes, bool, error) {
	// 判断该用户是否存在
	user := repositories.QueryUserById(u.DB, req.UserId)
	if user == nil {
		// return nil, fmt.Errorf("UserReqContext.Attention() : id为%d的用户不存在", req.UserId)
		globals.Log.Error(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(req.UserId)))
		return nil, false, fmt.Errorf(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(req.UserId)))
	}

	// 查询
	ids, isHaveData, err := repositories.QueryAttentionByPage(u.DB, req.UserId, req.Keyword, req.Page, req.Limit)
	if err != nil {
		// return nil, fmt.Errorf("UserReqContext.Attention() err: %v", err)
		globals.Log.Error(err.Error())
		return nil, false, err
	}

	// // 筛除数据（避免 ids 里面存在 user.Id ）
	// for k, v := range ids {
	// 	if v == req.UserId {
	// 		ids = append(ids[:k], ids[k+1:]...)
	// 	}
	// }

	res := &requests.AttentionRes{Ids: ids}

	return res, isHaveData, nil
}

// GetBasicInfo 通过ids获取到用户简略信息
func (u *UserReqContext) GetBasicInfo(req requests.GetBasicInfoReq) ([]*requests.GetBasicInfoRes, error) {
	userId := req.UserId
	// 存放结果
	var res = make([]*requests.GetBasicInfoRes, 0)

	// 游客登录
	if userId == 0 {

		for _, id := range req.Ids {
			// 查找用户id
			user := repositories.QueryUserById(u.DB, id)
			if user == nil {
				globals.Log.Error(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(id)))
				return nil, fmt.Errorf(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(id)))
			}

			// 查询用户的头像路径
			userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, id)
			if err != nil {
				globals.Log.Error(err.Error())
				return nil, err
			}
			// 没有图片
			if userImages == nil {
				globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(id)))
				return nil, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(id)))
			}
			avatarPath := (*userImages)[0]

			// 查询用户的文章数量
			authorArticles, err := repositories.QueryUserIDArticleNumOfPub(u.DB, id, "public")
			if err != nil {
				globals.Log.Error(err.Error())
				return nil, err
			}

			highlightName := internalUtils.Highlight(user.Nickname, req.Keyword)
			res = append(res, &requests.GetBasicInfoRes{
				ID:              user.ID,
				CreatedAt:       user.CreatedAt,
				UpdatedAt:       user.UpdatedAt,
				Nickname:        highlightName,
				Email:           user.Email,
				Heat:            user.Heat,
				AttentionCount:  user.AttentionCount,
				FansCount:       user.FansCount,
				PrivateSettings: user.PrivateSettings,
				Status:          user.Status,
				LastLoginTime:   user.LastLoginTime,
				IsFollowed:      0,
				AvatarPath:      avatarPath,
				AuthorArticles:  int(authorArticles),
			})
		}

	} else {

		// 查询 userId 全部关注的人
		ids, err := repositories.QueryAttention(u.DB, userId)
		if err != nil {
			// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() err: %v", err)
			globals.Log.Error(err.Error())
			return nil, err
		}

		for _, id := range req.Ids {
			// 查找用户id
			user := repositories.QueryUserById(u.DB, id)
			if user == nil {
				// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() : id为%d的用户不存在", userId)
				globals.Log.Error(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(id)))
				return nil, fmt.Errorf(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(id)))
			}

			// 查询用户的头像路径
			userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, id)
			if err != nil {
				// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() %v", err)
				globals.Log.Error(err.Error())
				return nil, err
			}
			// 没有图片
			if userImages == nil {
				// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() err = 无法找到id为%d的用户头像图片", userId)
				globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(id)))
				return nil, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(id)))
			}
			avatarPath := (*userImages)[0]

			// 判断 userID 是否关注 id。未关注：0，已关注：1。
			isFollowed := 0
			for _, v := range ids {
				if v == id {
					isFollowed = 1
				}
			}

			// 查询用户的文章数量
			authorArticles, err := repositories.QueryUserIDArticleNumOfPub(u.DB, id, "public")
			if err != nil {
				// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() err: %v", err)
				globals.Log.Error(err.Error())
				return nil, err
			}

			highlightName := internalUtils.Highlight(user.Nickname, req.Keyword)
			res = append(res, &requests.GetBasicInfoRes{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				// Nickname:        user.Nickname,
				Nickname:        highlightName,
				Email:           user.Email,
				Heat:            user.Heat,
				AttentionCount:  user.AttentionCount,
				FansCount:       user.FansCount,
				PrivateSettings: user.PrivateSettings,
				Status:          user.Status,
				LastLoginTime:   user.LastLoginTime,
				IsFollowed:      isFollowed,
				AvatarPath:      avatarPath,
				AuthorArticles:  int(authorArticles),
			})
		}

	}

	return res, nil

	// // 存放结果
	// var res = make([]*requests.GetBasicInfoRes, 0)
	//
	// // 查询 userId 全部关注的人
	// ids, err := repositories.QueryAttention(u.DB, userId)
	// if err != nil {
	// 	// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() err: %v", err)
	// 	globals.Log.Error(err.Error())
	// 	return nil, err
	// }
	//
	// for _, id := range req.Ids {
	// 	// 查找用户id
	// 	user := repositories.QueryUserById(u.DB, id)
	// 	if user == nil {
	// 		// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() : id为%d的用户不存在", userId)
	// 		globals.Log.Error(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(id)))
	// 		return nil, fmt.Errorf(response.ErrUserIdNotExist + ":" + strconv.Itoa(int(id)))
	// 	}
	//
	// 	// 查询用户的头像路径
	// 	userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, id)
	// 	if err != nil {
	// 		// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() %v", err)
	// 		globals.Log.Error(err.Error())
	// 		return nil, err
	// 	}
	// 	// 没有图片
	// 	if userImages == nil {
	// 		// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() err = 无法找到id为%d的用户头像图片", userId)
	// 		globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(id)))
	// 		return nil, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(id)))
	// 	}
	// 	avatarPath := (*userImages)[0]
	//
	// 	// 判断 userID 是否关注 id。未关注：0，已关注：1。
	// 	isFollowed := 0
	// 	for _, v := range ids {
	// 		if v == id {
	// 			isFollowed = 1
	// 		}
	// 	}
	//
	// 	// 查询用户的文章数量
	// 	authorArticles, err := repositories.QueryUserIDArticleNumOfPub(u.DB, id, "public")
	// 	if err != nil {
	// 		// return nil, fmt.Errorf("UserReqContext.GetBasicInfo() err: %v", err)
	// 		globals.Log.Error(err.Error())
	// 		return nil, err
	// 	}
	//
	// 	highlightName := internalUtils.Highlight(user.Nickname, req.Keyword)
	// 	res = append(res, &requests.GetBasicInfoRes{
	// 		ID:        user.ID,
	// 		CreatedAt: user.CreatedAt,
	// 		UpdatedAt: user.UpdatedAt,
	// 		// Nickname:        user.Nickname,
	// 		Nickname:        highlightName,
	// 		Email:           user.Email,
	// 		Heat:            user.Heat,
	// 		AttentionCount:  user.AttentionCount,
	// 		FansCount:       user.FansCount,
	// 		PrivateSettings: user.PrivateSettings,
	// 		Status:          user.Status,
	// 		LastLoginTime:   user.LastLoginTime,
	// 		IsFollowed:      isFollowed,
	// 		AvatarPath:      avatarPath,
	// 		AuthorArticles:  int(authorArticles),
	// 	})
	// }
	// return res, nil
}

// GetUserArticleLogic
// @Description: 获取用户的文章
// @param        db *gorm.DB
// @param        userId int
// @param        page int
// @param        limit int
// @return       []*models.Article
// @return       error
// @Author tianjiajie 2025-01-17 17:12:57
// func GetUserArticleLogic(db *gorm.DB, req requests.UserDataRequest, b bool) (data interface{}, err error) {
//	// 查询用户的文章
//	articles, total, err := repositories.QueryUserArticleRep(db, req, b)
//	if err != nil {
//		return nil, fmt.Errorf("GetUserArticleLogic() -> %v", err)
//	}
//
//	// 封装数据
//	data = gin.H{
//		"articles": articles,
//		"total":    total,
//	}
//
//	return data, nil
// }
