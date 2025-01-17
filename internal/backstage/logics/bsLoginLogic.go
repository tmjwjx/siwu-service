package logics

// BsManageContext
// @Description: 用于在处理请求时传递数据库连接和请求上下文信息
// @Author lizhuang 2024-10-21 20:26:10
// type BsManageContext struct {
// 	DB  *gorm.DB
// 	Ctx *gin.Context
// }

// NewBsManageContext
// @Description: 新建BsManageContext对象
// @Author lizhuang 2024-10-21 20:26:15
// @param        db *gorm.DB
// @param        c *gin.Context
// @return       *BsManageContext
// func NewBsManageContext(db *gorm.DB, c *gin.Context) *BsManageContext {
// 	return &BsManageContext{
// 		DB:  db,
// 		Ctx: c,
// 	}
// }
//
// // BsLogout 后台登出
// func (b *BsManageContext) BsLogout(tokenString string) error {
// 	// 使token无效
// 	if err := token.InvalidateToken(tokenString); err != nil {
// 		return fmt.Errorf("BsManageContext.BsLogout() err = %v", err)
// 	}
// 	return nil
// }
