package logics

import (
	"forum/internal/message/requests"
	"gorm.io/gorm"
)

func SystemMessageLogic(db *gorm.DB, req requests.MessageReq, id uint) (data interface{}, err error) {
	//res, err := repositories.CollectionRep(db, req, id)
	//data = gin.H{"collection_list": res}
	return data, err
}