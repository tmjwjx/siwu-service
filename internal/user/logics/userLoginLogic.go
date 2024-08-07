package logics

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
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
func (u *UserReqContext) Register(registerMsg *requests.RegisterMsg) error {
	var err error

	// 随机生成验证码
	// verifyCode := RandomGenerateStrings(VerifyCodeLen)

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
		err = repositories.InsertUser(u.DB, user)
		if err != nil {
			return fmt.Errorf("UserReqContext.Register() -> %s", err.Error())
		}
	}

	return nil
}

// ReqVerifyCode 用户请求验证码
func (u *UserReqContext) ReqVerifyCode(reqVerifyCode *requests.ReqVerifyCode) error {
	// 随机生成验证码
	verifyCode := internal_utils.RandomGenerateStrings(internal_utils.VerifyCodeLen)
	email := reqVerifyCode.Email

	// 给用户发送验证码
	body := fmt.Sprintf("你的验证码为 %s，有效时间为 %d 分钟\n", verifyCode, int(internal_utils.VerifyCodeEffectiveDuration.Minutes()))
	err := SendEmail(internal_utils.Form, email, internal_utils.Subject, body, internal_utils.AuthorizeCode)
	if err != nil {
		return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
	}
	//
	// // 存储数据
	//
	// // 有效时间
	// // 发送的冷却时间
	//
	// // 判断该email是否已经有用户使用过
	// user := repositories.QueryUserByEmail(u.DB, email)
	// // 如果已经有用户使用，而且发送验证码的冷却时间到了，更新UserVerifyCode表中的验证码
	// if user != nil {
	// 	userVerifyCode := repositories.QueryUserVerifyCodeByUID(u.DB, user.ID)
	// 	// 计算更新时间和当前时间的差异
	// 	duration := time.Now().Sub(userVerifyCode.UpdatedAt)
	// 	// 如果冷却时间未到，返回错误
	// 	if duration < internal_utils.VerifyCodeCoolTime {
	// 		return fmt.Errorf("UserReqContext.ReqVerifyCode() err: 发送验证码正在冷却时间中")
	// 	}
	//
	// 	if err = repositories.UpdateVerifyCodeByUID(u.DB, user.ID, verifyCode); err != nil {
	// 		return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
	// 	}
	//
	// } else { // 如果没有用户使用，向user表中插入用户，并向UserVerifyCode表中插入验证码
	// 	// 随机生成用户名
	// 	name := internal_utils.RandomGenerateStrings(8)
	// 	user := &models.User{
	// 		Nickname: name,
	// 		Email:    email,
	// 	}
	// 	// 向user表中插入新数据
	// 	if err = repositories.Insert(u.DB, user); err != nil {
	// 		return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
	// 	}
	//
	// 	// 将验证码插入到 UserVerifyCode表
	// 	var userVerifyCode = &models.UserVerifyCode{
	// 		UserID:     user.ID,
	// 		VerifyCode: verifyCode,
	// 	}
	// 	if err = repositories.Insert(u.DB, userVerifyCode); err != nil {
	// 		return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
	// 	}
	// }

	return nil
}
