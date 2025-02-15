package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FollowUnreadCountLogic
// @Description: 未读关注消息数量
// @param        db *gorm.DB
// @param        userId uint
// @return       data
// @return       err
// @Author tianjiajie 2025-02-12 22:12:33
func FollowUnreadCountLogic(db *gorm.DB, userId uint) (data interface{}, err error) {
	count, err := repositories.FollowUnreadCount(db, userId)
	if err != nil {
		return 0, err
	}
	data = gin.H{"count": count}
	return data, nil
}

// FollowMessageLogic
// @Description: 关注消息逻辑
// @Author tianjiajie 2025-01-16 20:29:36
func FollowMessageLogic(db *gorm.DB, req requests.MessageReq, userId uint) (data interface{}, err error) {
	res, err := repositories.FollowRep(db, req, userId)
	if err != nil {
		return nil, err
	}
	data = gin.H{"follow_list": res}
	return data, err
}
