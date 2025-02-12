package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CollectionUnreadCountLogic
// @Description: 未读收藏消息数量
// @param        db *gorm.DB
// @param        id uint
// @return       count
// @return       err
// @Author tianjiajie 2025-02-12 22:09:36
func CollectionUnreadCountLogic(db *gorm.DB, id uint) (data interface{}, err error) {
	count, err := repositories.CollectionUnreadCount(db, id)
	if err != nil {
		return 0, err
	}
	data = gin.H{"count": count}
	return data, err
}

// CollectionMessageLogic
// @Description: 收藏消息
// @param        db *gorm.DB
// @param        req requests.MessageReq
// @param        id uint
// @return       data
// @return       err
// @Author tianjiajie 2024-10-05 17:53:00
func CollectionMessageLogic(db *gorm.DB, req requests.MessageReq, id uint) (data interface{}, err error) {
	res, err := repositories.CollectionRep(db, req, id)
	data = gin.H{"collection_list": res}
	return data, err
}
