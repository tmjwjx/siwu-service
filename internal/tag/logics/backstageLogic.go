package logics

import (
	"fmt"
	"forum/internal/tag/repositories"
	"forum/internal/tag/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddTagLogic 新增标签
func AddTagLogic(c *gin.Context, db *gorm.DB, req *requests.BsAddTagReq) (error, int) {
	err, status := repositories.AddTagRep(c, db, req)
	return err, status
}

// DeleteTagLogic 删除标签
func DeleteTagLogic(db *gorm.DB, req *requests.BsDelTagReq) error {
	err := repositories.DeleteTagRep(db, req)
	return err
}

// BatchDelTagLogic 批量删除
func BatchDelTagLogic(db *gorm.DB, req *requests.BsBatchDelTagReq) error {
	err := repositories.BatchDelTagRep(db, req)
	return err
}

// UpdateTagLogic 更新标签
func UpdateTagLogic(c *gin.Context, db *gorm.DB, req *requests.BsUpTagReq) (error, int) {

	err, status := repositories.UpdateTagRep(c, db, req)

	return err, status
}

// QueryTagLogic 查询标签
func QueryTagLogic(db *gorm.DB, req *requests.BsQueTagReq) (*requests.BsQueTag, error) {
	tagRes, err := repositories.QueryTagRep(db, req)
	return tagRes, err
}

// BatchQueryTagLogic 批量查询标签
func BatchQueryTagLogic(db *gorm.DB, req *requests.BsBatchQueTagReq) (*requests.BsQueTagRes, error) {
	if req.Offset <= 0 || req.Limit <= 0 {
		return nil, fmt.Errorf("BatchQueryTagLogic -> Offset 或 Limit 的值不能小于或等于0")
	}
	batchTagRes, err := repositories.BatchQueryTagRep(db, req)
	return batchTagRes, err
}
