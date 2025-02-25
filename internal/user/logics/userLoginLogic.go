package logics

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/internalPkg/templates"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/flowRestriction"
	"forum/pkg/globals"
	"forum/pkg/response"
	"forum/pkg/sendEmailAsynchronous"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// UserReqContext 用于在处理请求时传递数据库连接和请求上下文信息
type UserReqContext struct {
	DB  *gorm.DB
	Ctx *gin.Context
}

// NewUserReqContext 新建UserReqContext对象
func NewUserReqContext(db *gorm.DB, c *gin.Context) *UserReqContext {
	return &UserReqContext{
		DB:  db,
		Ctx: c,
	}
}

// Register 注册
func (u *UserReqContext) Register(registerMsg requests.RegisterReq) (*requests.LogicRes, int, error) {
	email := registerMsg.Email
	verifyCode := registerMsg.VerifyCode
	password := registerMsg.Password

	// 判断是否已经注册过
	user := repositories.QueryUserByEmail(u.DB, email)
	if user != nil {
		globals.Log.Errorf(response.ErrEmailIsUse + ":" + email)
		return nil, 400, fmt.Errorf(response.ErrEmailIsUse + ":" + email)
		// return fmt.Errorf("UserReqContext.Register() : 邮箱为%s的用户已经注册过", email)
	}

	// 根据 email 查询该用户的最后一条验证码
	userVerifyCode, err := repositories.QueryLastUserVerifyCodeByEmail(u.DB, email)
	if err != nil {
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	// 判断该验证码是否使用过（判断DeletedAt是否有值）
	if userVerifyCode.DeletedAt.Valid {
		globals.Log.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
		return nil, 400, fmt.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
		// return fmt.Errorf("UserReqContext.Register() : 验证码%s已失效", verifyCode)
	}
	// 检验验证码是否正确（不区分大小写）
	if !strings.EqualFold(verifyCode, userVerifyCode.VerifyCode) {
		// return fmt.Errorf("UserReqContext.Register() : 验证码%s输入错误", verifyCode)
		globals.Log.Errorf(response.ErrVerifyCodeIsWrong + ":" + userVerifyCode.VerifyCode)
		return nil, 400, fmt.Errorf(response.ErrVerifyCodeIsWrong + ":" + userVerifyCode.VerifyCode)
	}
	// 判断验证码是否已经超时
	now := time.Now()
	duration := now.Sub(userVerifyCode.UpdatedAt)
	if duration > internalUtils.VerifyCodeEffectiveDuration {
		// return fmt.Errorf("UserReqContext.Register() : 验证码%s已过期", verifyCode)
		globals.Log.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
		return nil, 400, fmt.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
	}

	// 随机生成用户名
	// nickName := internalUtils.RandomGenerateStrings(internalUtils.UserNameLen)
	nickName, err := internalUtils.RandomGenerateNickname()
	if err != nil {
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	// 密码加密
	encryptedPassword, err := internalUtils.HashPassword(password)
	if err != nil {
		// return fmt.Errorf("UserReqContext.Register() : 密码%s加密失败", password)
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	// 向 user 表中添加该用户
	if err = sqlUtils.InsertObject(u.DB, &models.User{Nickname: nickName, Email: email, Password: encryptedPassword, LastLoginTime: time.Now()}); err != nil {
		// return fmt.Errorf("UserReqContext.Register() err: %v", err)
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}
	// 查询该用户的id
	user = repositories.QueryUserByEmail(u.DB, email)
	if user == nil {
		// return fmt.Errorf("UserReqContext.Register() : 未查询到邮箱为%v的用户", email)
		globals.Log.Errorf(response.ErrEmailNotExist + ":" + email)
		return nil, 500, fmt.Errorf(response.ErrEmailNotExist + ":" + email)
	}
	// 向 UserDetail 表中添加该用户
	if err = sqlUtils.InsertObject(u.DB, &models.UserDetail{UserID: user.ID}); err != nil {
		// return fmt.Errorf("UserReqContext.Register() err: %v", err)
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	// 删除该用户对应的全部验证码
	_, err = sqlUtils.DeleteObjectsByModel(u.DB, &models.UserVerifyCode{}, map[string]interface{}{"email": email})
	if err != nil {
		// return fmt.Errorf("UserReqContext.Register() err: %v", err)
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, user.ID)
	if err != nil {
		globals.Log.Errorf(err.Error())
		return nil, 500, err
		// return nil, fmt.Errorf("UserReqContext.Login() %v", err)
	}
	// 没有图片
	if userImages == nil {
		// return nil, fmt.Errorf("UserReqContext.Login() err = 无法找到id为%d的用户头像图片", user.ID)
		globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(user.ID)))
		return nil, 500, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(user.ID)))
	}
	avatarPath := (*userImages)[0]

	/*// 为用户设置默认头像
	da := internalUtils.NewDefaultAvatar(globals.UserHome, user.ID, u.DB)
	avatarPath, err := internalUtils.GenerateAvatar(da)
	if err != nil {
		globals.Log.Errorf(err.Error())
	}*/

	// 修改 LastLoginTime
	if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: user.ID}}, map[string]interface{}{"last_login_time": now}); err != nil {
		// return nil, fmt.Errorf("UserReqContext.Login() -> %v", err)
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	// 获取登陆响应
	userInfo := &requests.LogicRes{
		Id:         user.ID,
		Nickname:   user.Nickname,
		AvatarPath: avatarPath,
	}

	return userInfo, 200, nil
}

// ReqVerifyCode 用户请求验证码
func (u *UserReqContext) ReqVerifyCode(email string) (int, error) {
	// 查找最后一条验证码（如何能够查询到验证码，要判断一下冷却时间；如果不能够查询到验证码，就直接发送）
	userVerifyCode, _ := repositories.QueryLastUserVerifyCodeByEmail(u.DB, email)
	// 如果查询到了验证码，判断冷却时间是否到
	if userVerifyCode != nil {
		// 当前时间
		now := time.Now()
		// 计算更新时间和当前时间的差异
		duration := now.Sub(userVerifyCode.UpdatedAt)
		// 如果冷却时间未到，返回错误
		if duration < internalUtils.VerifyCodeCoolTime {
			// return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 发送验证码正在冷却时间中")
			globals.Log.Errorf(response.ErrReqVerifyCodeIsCooling)
			return 400, fmt.Errorf(response.ErrReqVerifyCodeIsCooling)
		}
	}

	// 随机生成验证码
	// verifyCode := internalUtils.RandomGenerateStrings(internalUtils.VerifyCodeLen)
	verifyCode, err := internalUtils.RandomGenerateVerifyCode()
	if err != nil {
		globals.Log.Errorf(err.Error())
		return 500, err
	}

	// 先发送验证码，后插入到表中，避免发送验证码失败
	// 给用户发送验证码
	// body := fmt.Sprintf(templates.GetEmailFormatTemplate(), verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
	// if err := globals.SendEmailCfg.SendEmail(email, internalUtils.VerifyCodeSubject, body); err != nil {
	// 	// return fmt.Errorf("UserReqContext.VerifyCodeReq() -> 向 %s 邮箱发送验证码错误，%v", email, err)
	// 	globals.Log.Errorf(err.Error())
	// 	return err
	// }

	body := fmt.Sprintf(templates.GetEmailFormatTemplate(), verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
	// if err := SendEmail(globals.SendEmailCfg, email, internalUtils.VerifyCodeSubject, body); err != nil {
	// 	// return fmt.Errorf("UserReqContext.VerifyCodeReq() -> 向 %s 邮箱发送验证码错误，%v", email, err)
	// 	globals.Log.Errorf(err.Error())
	// 	return err
	// }

	// 创建邮件任务消息
	task := sendEmailAsynchronous.EmailTask{
		Email:   email,
		Subject: internalUtils.VerifyCodeSubject,
		Body:    body,
	}

	// 将任务推送到 Redis Stream（异步处理）
	if err := sendEmailAsynchronous.PushEmailTaskToStream(globals.RDB, task); err != nil {
		globals.Log.Errorf("推送邮件任务到 Redis Stream 失败: %v", err)
		return 500, err
	}

	// 插入一条新的验证码数据
	if err := sqlUtils.InsertObject(u.DB, &models.UserVerifyCode{Email: email, VerifyCode: verifyCode}); err != nil {
		// return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
		globals.Log.Errorf(err.Error())
		return 500, err
	}

	return 200, nil
}

// Login 登录
func (u *UserReqContext) Login(logicMsg requests.LogicReq) (*requests.LogicRes, int, error) {
	// 判断邮箱和密码是否匹配
	email := logicMsg.Email
	password := logicMsg.Password

	// 根据邮箱查用户
	user := repositories.QueryUserByEmail(u.DB, email)
	if user == nil {
		// return nil, fmt.Errorf("UserReqContext.Login() err: 不存在该邮箱用户")
		globals.Log.Errorf(response.ErrEmailNotExist + ":" + email)
		return nil, 400, fmt.Errorf(response.ErrEmailNotExist + ":" + email)
	}

	userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, user.ID)
	if err != nil {
		globals.Log.Errorf(err.Error())
		return nil, 500, err
		// return nil, fmt.Errorf("UserReqContext.Login() %v", err)
	}
	// 没有图片
	if userImages == nil {
		// return nil, fmt.Errorf("UserReqContext.Login() err = 无法找到id为%d的用户头像图片", user.ID)
		globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(user.ID)))
		return nil, 500, fmt.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(user.ID)))
	}
	avatarPath := (*userImages)[0]

	// 比较加密密码
	encryptedPassword := user.Password
	if !internalUtils.CheckPasswordHash(password, encryptedPassword) {
		// return nil, fmt.Errorf("UserReqContext.Login() err: 密码错误")

		// 记录登录失败（限流）
		flowRestriction.RecordFailedAttempt(globals.RDB, u.Ctx)

		globals.Log.Errorf(response.ErrPasswordIsWrong)
		return nil, 400, fmt.Errorf(response.ErrPasswordIsWrong)
	}

	// 修改 LastLoginTime
	now := time.Now() // 获取当前时间
	if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: user.ID}}, map[string]interface{}{"last_login_time": now}); err != nil {
		// return nil, fmt.Errorf("UserReqContext.Login() -> %v", err)
		globals.Log.Errorf(err.Error())
		return nil, 500, err
	}

	// 获取登陆响应
	logicRes := &requests.LogicRes{
		Id:         user.ID,
		Nickname:   user.Nickname,
		AvatarPath: avatarPath,
	}
	return logicRes, 200, nil
}

// ForgotPassword 忘记密码
func (u *UserReqContext) ForgotPassword(forgotPasswordMsg requests.ForgotPasswordReq) (int, error) {
	email := forgotPasswordMsg.Email
	verifyCode := forgotPasswordMsg.VerifyCode
	password := forgotPasswordMsg.Password

	// 判断是否已经注册过
	user := repositories.QueryUserByEmail(u.DB, email)
	// 必须已经注册过该用户
	if user == nil {
		// return fmt.Errorf("UserReqContext.Register() : 邮箱为%s的用户没有注册过", email)
		globals.Log.Errorf(response.ErrEmailNotExist + ":" + email)
		return 400, fmt.Errorf(response.ErrEmailNotExist + ":" + email)
	}

	// 根据 email 查询该用户的最后一条验证码
	userVerifyCode, err := repositories.QueryLastUserVerifyCodeByEmail(u.DB, email)
	if err != nil {
		// return fmt.Errorf("UserReqContext.Register() -> %v", err)
		globals.Log.Errorf(err.Error())
		return 500, err
	}

	// 判断该验证码是否使用过（判断DeletedAt是否有值）
	if userVerifyCode.DeletedAt.Valid {
		// return fmt.Errorf("UserReqContext.Register() : 验证码%s已失效", verifyCode)
		globals.Log.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
		return 400, fmt.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
	}
	// 检验验证码是否正确（不区分大小写）
	if !strings.EqualFold(verifyCode, userVerifyCode.VerifyCode) {
		// return fmt.Errorf("UserReqContext.Register() : 验证码%s输入错误", verifyCode)
		globals.Log.Errorf(response.ErrVerifyCodeIsWrong + ":" + verifyCode)
		return 400, fmt.Errorf(response.ErrVerifyCodeIsWrong + ":" + verifyCode)
	}
	// 判断验证码是否已经超时
	now := time.Now()
	duration := now.Sub(userVerifyCode.UpdatedAt)
	if duration > internalUtils.VerifyCodeEffectiveDuration {
		// return fmt.Errorf("UserReqContext.Register() : 验证码%s已过期", verifyCode)
		globals.Log.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
		return 400, fmt.Errorf(response.ErrVerifyCodeIsExpired + ":" + userVerifyCode.VerifyCode)
	}

	// 密码加密
	encryptedPassword, err := internalUtils.HashPassword(password)
	if err != nil {
		// return fmt.Errorf("UserReqContext.Register() : 密码%s加密失败", password)
		globals.Log.Errorf(err.Error())
		return 500, err
	}

	// 更新密码
	if err = sqlUtils.UpdateObjects(u.DB, &models.User{Email: email}, map[string]interface{}{"password": encryptedPassword}); err != nil {
		// return fmt.Errorf("UserReqContext.Register() err: %v", err)
		globals.Log.Errorf(err.Error())
		return 500, err
	}

	// 删除该用户对应的全部验证码
	_, err = sqlUtils.DeleteObjectsByModel(u.DB, &models.UserVerifyCode{}, map[string]interface{}{"email": email})
	if err != nil {
		// return fmt.Errorf("UserReqContext.Register() err: %v", err)
		globals.Log.Errorf(err.Error())
		return 500, err
	}

	return 200, nil
}

// Logout 登出
func (u *UserReqContext) Logout(tokenString string) error {
	// 设置过期时间为 Token 剩余时间
	claims, err := token.ValidateToken(tokenString)
	if err != nil {
		globals.Log.Errorf(response.ErrTokenIsInvalid)
		return fmt.Errorf(response.ErrTokenIsInvalid)
	}

	expiration := time.Until(claims.ExpiresAt.Time)
	// 如果该token还没有失效，就把它添加到黑名单中，让它失效
	if expiration > 0 {
		if err = token.AddTokenToBlacklist(globals.RDB, tokenString, expiration); err != nil {
			globals.Log.Errorf(err.Error())
			return err
		}
	}
	return nil
}
