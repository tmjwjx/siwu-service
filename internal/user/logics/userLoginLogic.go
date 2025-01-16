package logics

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/internalPkg/templates"
	"strings"

	// "forum/internal/internalPkg/templates"
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
	DB           *gorm.DB
	Ctx          *gin.Context
	SendEmailCfg *globals.SendEmailConfig // 发送邮件
}

// NewUserReqContext
// @Description: 新建UserReqContext对象
// @Author lizhuang 2024-10-09 09:38:14
// @param        db *gorm.DB
// @param        c *gin.Context
// @param        sendEmailCfg *globals.SendEmailConfig
// @return       *UserReqContext
func NewUserReqContext(db *gorm.DB, c *gin.Context, sendEmailCfg *globals.SendEmailConfig) *UserReqContext {
	return &UserReqContext{
		DB:           db,
		Ctx:          c,
		SendEmailCfg: sendEmailCfg,
	}
}

// Register 注册
func (u *UserReqContext) Register(registerMsg requests.RegisterReq) error {
	email := registerMsg.Email
	verifyCode := registerMsg.VerifyCode
	password := registerMsg.Password

	// 判断是否已经注册过
	user := repositories.QueryUserByEmail(u.DB, email)
	if user != nil {
		return fmt.Errorf("UserReqContext.Register() : 邮箱为%s的用户已经注册过", email)
	}

	// 根据 email 查询该用户的最后一条验证码
	userVerifyCode, err := repositories.QueryLastUserVerifyCodeByEmail(u.DB, email)
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() -> %v", err)
	}

	// 判断该验证码是否使用过（判断DeletedAt是否有值）
	if userVerifyCode.DeletedAt.Valid {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s已失效", verifyCode)
	}
	// 检验验证码是否正确（不区分大小写）
	if !strings.EqualFold(verifyCode, userVerifyCode.VerifyCode) {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s输入错误", verifyCode)
	}
	// 判断验证码是否已经超时
	now := time.Now()
	duration := now.Sub(userVerifyCode.UpdatedAt)
	if duration > internalUtils.VerifyCodeEffectiveDuration {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s已过期", verifyCode)
	}

	// 随机生成用户名
	nickName := internalUtils.RandomGenerateStrings(internalUtils.UserNameLen)

	// 密码加密
	encryptedPassword, err := internalUtils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() : 密码%s加密失败", password)
	}

	// 向 user 表中添加该用户
	if err = sqlUtils.InsertObject(u.DB, &models.User{Nickname: nickName, Email: email, Password: encryptedPassword, LastLoginTime: time.Time{}}); err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}

	// 删除该用户对应的全部验证码
	_, err = sqlUtils.DeleteObjectsByModel(u.DB, &models.UserVerifyCode{}, map[string]interface{}{"email": email})
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}
	return nil
}

// ReqVerifyCode 用户请求验证码
func (u *UserReqContext) ReqVerifyCode(email string) error {
	// 随机生成验证码
	verifyCode := internalUtils.RandomGenerateStrings(internalUtils.VerifyCodeLen)

	// 查找最后一条验证码（如何能够查询到验证码，要判断一下冷却时间；如果不能够查询到验证码，就直接发送）
	userVerifyCode, _ := repositories.QueryLastUserVerifyCodeByEmail(u.DB, email)
	// 如何查询到了验证码，判断冷却时间是否到
	if userVerifyCode != nil {
		// 当前时间
		now := time.Now()
		// 计算更新时间和当前时间的差异
		duration := now.Sub(userVerifyCode.UpdatedAt)
		// 如果冷却时间未到，返回错误
		if duration < internalUtils.VerifyCodeCoolTime {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 发送验证码正在冷却时间中")
		}
	}

	// 插入一条新的验证码数据
	if err := sqlUtils.InsertObject(u.DB, &models.UserVerifyCode{Email: email, VerifyCode: verifyCode}); err != nil {
		return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
	}

	// 给用户发送验证码
	body := fmt.Sprintf(templates.GetEmailFormatTemplate(), verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
	if err := u.SendEmail(email, internalUtils.VerifyCodeSubject, body); err != nil {
		return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
	}

	return nil
}

// SendEmail 给邮箱(to)发送内容(body)
// to: 接收人
// body: 正文内容
// subject: 主题
func (u *UserReqContext) SendEmail(to string, subject string, body string) error {
	// 判断邮箱是否合法
	if !internalUtils.IsValidEmail(to) {
		return fmt.Errorf("UserReqContext.SendEmail() err: 接收者邮箱错误")
	}

	m := gomail.NewMessage()
	// 设置邮件消息的头部字段
	m.SetHeader("From", u.SendEmailCfg.From) // 发送人
	m.SetHeader("To", to)                    // 接收人
	m.SetHeader("Subject", subject)          // 主题
	m.SetBody("text/html", body)             // 正文内容
	// 创建一个新的邮件拨号器对象，用于通过指定的 SMTP 服务器发送邮件
	d := gomail.NewDialer(u.SendEmailCfg.Host, u.SendEmailCfg.Port, u.SendEmailCfg.Username, u.SendEmailCfg.AuthorizeCode)
	// 通过拨号器对象发送指定的邮件消息
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("UserReqContext.SendEmail() err: %v", err)
	}

	return nil
}

// Login 登录
func (u *UserReqContext) Login(logicMsg requests.LogicReq) (*requests.LogicRes, error) {
	// 判断邮箱和密码是否匹配
	email := logicMsg.Email
	password := logicMsg.Password

	// 根据邮箱查用户
	user := repositories.QueryUserByEmail(u.DB, email)
	if user == nil {
		return nil, fmt.Errorf("UserReqContext.Login() err: 不存在该邮箱用户")
	}

	// 比较加密密码
	encryptedPassword := user.Password
	if !internalUtils.CheckPasswordHash(password, encryptedPassword) {
		return nil, fmt.Errorf("UserReqContext.Login() err: 密码错误")
	}

	// 修改 LastLoginTime
	now := time.Now() // 获取当前时间
	if err := sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: user.ID}}, map[string]interface{}{"last_login_time": now}); err != nil {
		return nil, fmt.Errorf("UserReqContext.Login() -> %v", err)
	}

	// 获取登陆响应
	logicRes := &requests.LogicRes{
		Id:       user.ID,
		Nickname: user.Nickname,
	}
	return logicRes, nil
}

// ForgotPassword 忘记密码
func (u *UserReqContext) ForgotPassword(forgotPasswordMsg requests.ForgotPasswordReq) error {
	email := forgotPasswordMsg.Email
	verifyCode := forgotPasswordMsg.VerifyCode
	password := forgotPasswordMsg.Password

	// 判断是否已经注册过
	user := repositories.QueryUserByEmail(u.DB, email)
	// 必须已经注册过该用户
	if user == nil {
		return fmt.Errorf("UserReqContext.Register() : 邮箱为%s的用户没有注册过", email)
	}

	// 根据 email 查询该用户的最后一条验证码
	userVerifyCode, err := repositories.QueryLastUserVerifyCodeByEmail(u.DB, email)
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() -> %v", err)
	}

	// 判断该验证码是否使用过（判断DeletedAt是否有值）
	if userVerifyCode.DeletedAt.Valid {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s已失效", verifyCode)
	}
	// 检验验证码是否正确（不区分大小写）
	if !strings.EqualFold(verifyCode, userVerifyCode.VerifyCode) {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s输入错误", verifyCode)
	}
	// 判断验证码是否已经超时
	now := time.Now()
	duration := now.Sub(userVerifyCode.UpdatedAt)
	if duration > internalUtils.VerifyCodeEffectiveDuration {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s已过期", verifyCode)
	}

	// 密码加密
	encryptedPassword, err := internalUtils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() : 密码%s加密失败", password)
	}

	// 更新密码
	if err = sqlUtils.UpdateObjects(u.DB, &models.User{Email: email}, map[string]interface{}{"password": encryptedPassword}); err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}

	// 删除该用户对应的全部验证码
	_, err = sqlUtils.DeleteObjectsByModel(u.DB, &models.UserVerifyCode{}, map[string]interface{}{"email": email})
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}
	return nil
}
