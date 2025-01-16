package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
