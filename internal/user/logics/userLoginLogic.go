package logics

import (
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserReqContext 用于在处理请求时传递数据库连接和请求上下文信息
type UserReqContext struct {
	DB  *gorm.DB
	Ctx *gin.Context
}

func NewUserLogic(db *gorm.DB, c *gin.Context) *UserReqContext {
	return &UserReqContext{
		DB:  db,
		Ctx: c,
	}
}

// Register 注册
func (u *UserReqContext) Register(registerMsg requests.RegisterMsg) error {
	var err error

	// 随机生成验证码
	// verifyCode := RandomVerifyCode(VerifyCodeLen)

	// 检验邮箱是否合理，检验密码是否合理

	// 核对两次输入的密码

	// 判断邮箱是否已经注册过；

	// 验证码核对

	// 该邮箱未被注册
	// 获取验证码
	// 存储验证码
	// 发送验证码

	// 添加到数据库中
	if true {
		user := &models.User{
			Email: registerMsg.Email,
			// 密码要加密
			// todo
			Password: registerMsg.Password,
		}
		err = repositories.CreateUser(u.DB, user)
		if err != nil {
			return err
		}
	}

	return nil
}

// // ReqVerifyCode 用户请求验证码
// func (u *UserReqContext) ReqVerifyCode() (string, error) {
//
// 	return
// }
