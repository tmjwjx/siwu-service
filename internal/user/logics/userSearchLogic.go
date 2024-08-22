package logics

import (
	"fmt"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/globals"
)

// UserRank 用户热度排行
func (u *UserReqContext) UserRank(rankMsg requests.RankMsg) ([]*requests.UserDataRes, error) {
	limit := rankMsg.Limit

	// 查询用户ID排行
	userIdSli, err := repositories.QueryUserRank(u.DB, limit)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.ReqUserMsg() -> %v", err)
	}

	var userResponses = make([]*requests.UserDataRes, 0)
	// 根据ID切片查询各个用户的详细信息，如昵称，头像，个人简介
	for _, userID := range userIdSli {
		userResponse, err := repositories.UserDataResponse(userID, globals.DB)
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.ReqUserMsg() -> %v", err)
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}
