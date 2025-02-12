package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LikeUnreadCountLogic
// @Description: 未读点赞消息数量
// @param        db *gorm.DB
// @param        id uint
// @return       count
// @return       err
// @Author tianjiajie 2025-02-12 21:30:52
func LikeUnreadCountLogic(db *gorm.DB, id uint) (data interface{}, err error) {
	count, err := repositories.LikeUnreadCount(db, id)
	if err != nil {
		return 0, err
	}
	data = gin.H{"count": count}
	return data, nil
}

// LikeMessageLogic
// @Description: 点赞消息
// @param        db *gorm.DB
// @param        req requests.MessageReq
// @param        id uint
// @return       data
// @return       err
// @Author tianjiajie 2024-10-04 20:54:53
func LikeMessageLogic(db *gorm.DB, req requests.MessageReq, id uint) (data interface{}, err error) {
	res, err := repositories.LikeRep(db, req, id)
	if err != nil {
		return nil, err
	}
	data = gin.H{"like_list": res}
	return data, nil
}
