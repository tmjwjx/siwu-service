package logics

import (
	"fmt"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"strconv"
)

// UserRank 用户热度排行
func (u *UserReqContext) UserRank(rankMsg requests.RankMsg) ([]*requests.UserResponse, error) {
	limit := rankMsg.Limit

	// 查询用户ID排行
	userIdSli, err := repositories.QueryUserRank(u.DB, limit)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.ReqUserMsg() -> %v", err)
	}

	var userResponses = make([]*requests.UserResponse, 0)
	// 根据ID切片查询各个用户的详细信息，如昵称，头像，个人简介
	for _, userID := range userIdSli {
		userResponse, err := repositories.SelectPersonData(strconv.Itoa(int(userID)))
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.ReqUserMsg() -> %v", err)
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}
