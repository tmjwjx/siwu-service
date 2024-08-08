package logics

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"time"
)

// UserReqContext 用于在处理请求时传递数据库连接和请求上下文信息
type UserReqContext struct {
	DB         *gorm.DB
	Ctx        *gin.Context
	VerifyCode *globals.VerifyCodeConfig // 发送邮件
}

func NewUserLogic(db *gorm.DB, c *gin.Context, verifyCode *globals.VerifyCodeConfig) *UserReqContext {
	return &UserReqContext{
		DB:         db,
		Ctx:        c,
		VerifyCode: verifyCode,
	}
}

// Register 注册
func (u *UserReqContext) Register(registerMsg *requests.RegisterMsg) error {
	// 在数据库中完善数据（用户获取验证码时已插入数据）

	// 验证码有效时间
	// 只能使用一次

	// 密码要加密
	// todo

	// if err := repositories.Updat(u.DB); err != nil {
	// 	return fmt.Errorf("UserReqContext.Register() err: %s", err.Error())
	// }

	return nil
}

// ReqVerifyCode 用户请求验证码
func (u *UserReqContext) ReqVerifyCode(reqVerifyCode *requests.ReqVerifyCode) error {
	// 随机生成验证码
	verifyCode := internal_utils.RandomGenerateStrings(internal_utils.VerifyCodeLen)
	email := reqVerifyCode.Email

	// 存储数据

	// 判断该email是否已经有用户使用过
	user := repositories.QueryUserByEmail(u.DB, email)
	// 如果已经有用户使用，而且发送验证码的冷却时间到了，更新UserVerifyCode表中的验证码
	if user != nil {
		userVerifyCode := repositories.QueryUserVerifyCodeByUID(u.DB, user.ID)
		// 当前时间
		now := time.Now()
		// 计算更新时间和当前时间的差异
		duration := now.Sub(userVerifyCode.UpdatedAt)

		fmt.Println(now, userVerifyCode.UpdatedAt, time.Now(), duration)

		// 如果冷却时间未到，返回错误
		if duration < internal_utils.VerifyCodeCoolTime {
			return fmt.Errorf("UserReqContext.ReqVerifyCode() err: 发送验证码正在冷却时间中")
		}

		if err := repositories.UpdateVerifyCodeByUID(u.DB, user.ID, verifyCode); err != nil {
			return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
		}

	} else { // 如果没有用户使用，向user表中插入用户，并向UserVerifyCode表中插入验证码
		// 随机生成用户名
		name := internal_utils.RandomGenerateStrings(8)
		user := &models.User{
			Nickname: name,
			Email:    email,
		}

		// 向user表中插入新数据
		if err := repositories.Insert(u.DB, user); err != nil {
			return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
		}

		// 将验证码插入到 UserVerifyCode表
		var userVerifyCode = &models.UserVerifyCode{
			UserID:     user.ID,
			VerifyCode: verifyCode,
		}
		if err := repositories.Insert(u.DB, userVerifyCode); err != nil {
			return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
		}
	}

	// 给用户发送验证码
	body := fmt.Sprintf("你的验证码为 %s，有效时间为 %d 分钟\n", verifyCode, int(internal_utils.VerifyCodeEffectiveDuration.Minutes()))
	err := u.SendEmail(email, body)
	if err != nil {
		return fmt.Errorf("UserReqContext.ReqVerifyCode() -> %s", err.Error())
	}

	return nil
}

// SendEmail 给邮箱(to)发送内容(body)
func (u *UserReqContext) SendEmail(to string, body string) error {
	// // 判断邮箱是否合法
	// if !internal_utils.IsValidEmail(to) {
	// 	return fmt.Errorf("UserReqContext.SendEmail() err: 接收者邮箱错误")
	// }

	m := gomail.NewMessage()
	// 设置邮件消息的头部字段
	m.SetHeader("From", u.VerifyCode.From)       // 发送人
	m.SetHeader("To", to)                        // 接收人
	m.SetHeader("Subject", u.VerifyCode.Subject) // 主题
	m.SetBody("text/plain", body)                // 正文内容
	// 创建一个新的邮件拨号器对象，用于通过指定的 SMTP 服务器发送邮件
	d := gomail.NewDialer(u.VerifyCode.Host, u.VerifyCode.Port, u.VerifyCode.Username, u.VerifyCode.AuthorizeCode)
	// 通过拨号器对象发送指定的邮件消息
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("UserReqContext.SendEmail() err: %s", err.Error())
	}

	return nil
}
