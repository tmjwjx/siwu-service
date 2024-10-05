package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
	data = gin.H{"like_list": res}
	return data, err
}