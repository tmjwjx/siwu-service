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
	"strings"
	"time"
)

// UserReqContext 用于在处理请求时传递数据库连接和请求上下文信息
type UserReqContext struct {
	DB           *gorm.DB
	Ctx          *gin.Context
	SendEmailCfg *globals.SendEmailConfig // 发送邮件
}

func NewUserLogic(db *gorm.DB, c *gin.Context, sendEmailCfg *globals.SendEmailConfig) *UserReqContext {
	return &UserReqContext{
		DB:           db,
		Ctx:          c,
		SendEmailCfg: sendEmailCfg,
	}
}

// Register 注册
func (u *UserReqContext) Register(registerMsg requests.RegisterMsg) error {
	// 在数据库中完善数据（用户获取验证码时已在数据库中创建了User，UserVerifyCode数据）
	email := registerMsg.Email
	verifyCode := registerMsg.VerifyCode
	password := registerMsg.Password
	user := repositories.QueryUserByEmail(u.DB, email)
	if user == nil {
		return fmt.Errorf("UserReqContext.Register() : 不存在该邮箱为%s的用户", email)
	}
	userID := user.ID

	// 验证码核对（不区分大小写）
	userVerifyCode, err := repositories.QueryLastUserVerifyCodeByUserID(u.DB, userID)
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() -> %v", err)
	}
	// 判断该验证码是否使用过（判断DeletedAt是否有值）
	if userVerifyCode.DeletedAt.Valid {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s已失效", verifyCode)
	}
	// 判断用户输入的验证码是否等于数据库中的验证码
	if !strings.EqualFold(verifyCode, userVerifyCode.VerifyCode) {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s输入错误", verifyCode)
	}
	// 判断验证码是否已经超时
	now := time.Now()
	duration := now.Sub(userVerifyCode.UpdatedAt)
	if duration > internal_utils.VerifyCodeEffectiveDuration {
		return fmt.Errorf("UserReqContext.Register() : 验证码%s已过期", verifyCode)
	}

	// 密码加密
	encryptedPassword, err := internal_utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() : 密码%s加密失败", password)
	}

	// 更新用户密码
	// err = repositories.UpdateObjects(u.DB, &models.User{}, map[string]interface{}{"id": userID, "email": email}, map[string]interface{}{"password": encryptedPassword})
	err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: userID}, Email: email}, map[string]interface{}{"password": encryptedPassword})
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}

	// 删除该用户对应的全部验证码
	_, err = repositories.DeleteObjectsByModel(u.DB, &models.UserVerifyCode{}, map[string]interface{}{"user_id": userID})
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}
	return nil
}

// ReqVerifyCode 用户请求验证码
func (u *UserReqContext) ReqVerifyCode(reqVerifyCode requests.VerifyCodeMsg) error {
	// 随机生成验证码
	verifyCode := internal_utils.RandomGenerateStrings(internal_utils.VerifyCodeLen)
	email := reqVerifyCode.Email

	// 存储数据

	// 判断该email是否已经有用户使用过
	user := repositories.QueryUserByEmail(u.DB, email)
	// 如果没有用户使用过这个email，向user表中插入用户，并向UserVerifyCode表中插入验证码
	if user == nil {
		// 随机生成用户名
		name := internal_utils.RandomGenerateStrings(internal_utils.UserNameLen)
		// 给用户生成一个默认密码
		password := internal_utils.RandomGenerateStrings(12)
		// 使用InsertObject()方法向user表中插入新数据，model参数必须是指针类型
		if err := repositories.InsertObject(u.DB, &models.User{Nickname: name, Email: email, Password: password}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeMsg() -> %v", err)
		}

		// 获取用户id
		user = repositories.QueryUserByEmail(u.DB, email)
		// 将验证码插入到 UserVerifyCode表
		if err := repositories.InsertObject(u.DB, &models.UserVerifyCode{UserID: user.ID, VerifyCode: verifyCode}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeMsg() -> %v", err)
		}

	} else { // 如果已经有用户使用，而且发送验证码的冷却时间到了，插入一条数据
		userVerifyCode, err := repositories.QueryLastUserVerifyCodeByUserID(u.DB, user.ID)
		if err != nil { // 执行错误，没有查询到验证码（可能是手动删除了数据库中的验证码，所以报错）
			// 插入一条新的验证码数据
			if err = repositories.InsertObject(u.DB, &models.UserVerifyCode{UserID: user.ID, VerifyCode: verifyCode}); err != nil {
				return fmt.Errorf("UserReqContext.VerifyCodeMsg() -> %v", err)
			}
			// 给用户发送验证码
			body := fmt.Sprintf("你的验证码为 %s，有效时间为 %d 分钟\n", verifyCode, int(internal_utils.VerifyCodeEffectiveDuration.Minutes()))
			err := u.SendEmail(email, internal_utils.VerifyCodeSubject, body)
			if err != nil {
				return fmt.Errorf("UserReqContext.VerifyCodeMsg() -> %v", err)
			}
			return nil
		}

		// 判断冷却时间
		// 当前时间
		now := time.Now()
		// 计算更新时间和当前时间的差异
		duration := now.Sub(userVerifyCode.UpdatedAt)
		// 如果冷却时间未到，返回错误
		if duration < internal_utils.VerifyCodeCoolTime {
			return fmt.Errorf("UserReqContext.VerifyCodeMsg() err: 发送验证码正在冷却时间中")
		}

		// 插入一条新的验证码数据
		if err = repositories.InsertObject(u.DB, &models.UserVerifyCode{UserID: user.ID, VerifyCode: verifyCode}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeMsg() -> %v", err)
		}
	}

	// 给用户发送验证码
	body := fmt.Sprintf("你的验证码为 %s，有效时间为 %d 分钟\n", verifyCode, int(internal_utils.VerifyCodeEffectiveDuration.Minutes()))
	err := u.SendEmail(email, internal_utils.VerifyCodeSubject, body)
	if err != nil {
		return fmt.Errorf("UserReqContext.VerifyCodeMsg() -> %v", err)
	}

	return nil
}

// SendEmail 给邮箱(to)发送内容(body)
// to: 接收人
// body: 正文内容
// subject: 主题
func (u *UserReqContext) SendEmail(to string, subject string, body string) error {
	// 判断邮箱是否合法
	if !internal_utils.IsValidEmail(to) {
		return fmt.Errorf("UserReqContext.SendEmail() err: 接收者邮箱错误")
	}

	m := gomail.NewMessage()
	// 设置邮件消息的头部字段
	m.SetHeader("From", u.SendEmailCfg.From) // 发送人
	m.SetHeader("To", to)                    // 接收人
	m.SetHeader("Subject", subject)          // 主题
	m.SetBody("text/plain", body)            // 正文内容
	// 创建一个新的邮件拨号器对象，用于通过指定的 SMTP 服务器发送邮件
	d := gomail.NewDialer(u.SendEmailCfg.Host, u.SendEmailCfg.Port, u.SendEmailCfg.Username, u.SendEmailCfg.AuthorizeCode)
	// 通过拨号器对象发送指定的邮件消息
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("UserReqContext.SendEmail() err: %v", err)
	}

	return nil
}

// Login 登录
func (u *UserReqContext) Login(logicMsg requests.LogicMsg) error {
	// 判断邮箱和密码是否匹配
	email := logicMsg.Email
	password := logicMsg.Password

	// 根据邮箱查用户
	user := repositories.QueryUserByEmail(u.DB, email)
	if user == nil {
		return fmt.Errorf("UserReqContext.Login() : 不存在该邮箱用户")
	}

	// 比较加密密码
	encryptedPassword := user.Password
	if !internal_utils.CheckPasswordHash(password, encryptedPassword) {
		return fmt.Errorf("UserReqContext.Login() : 密码错误")
	}

	return nil
}

// Follow 关注。followerId 关注 followedId
func (u *UserReqContext) Follow(follow requests.FollowMsg) error {
	followerId := follow.FollowerId
	followedId := follow.FollowedId

	// 查询两个id，判断两个id是否存在
	// 关注者
	user1 := repositories.QueryUserById(u.DB, followerId)
	if user1 == nil {
		return fmt.Errorf("UserReqContext.Follow() : 关注者不存在")
	}
	// 被关注者
	user2 := repositories.QueryUserById(u.DB, followedId)
	if user2 == nil {
		return fmt.Errorf("UserReqContext.Follow() : 被关注者不存在")
	}

	// 判断是否已经关注过了，如果已经关注过了，再次点击就会取消关注
	// 我的关注
	followedIDSli, err := repositories.QueryFollowed(u.DB, followerId)
	if err != nil {
		return fmt.Errorf("UserReqContext.Follow() -> %v: ", err)
	}
	// 是否已经关注过 followedId
	var isFollowed = false
	for _, id := range followedIDSli {
		if id == followedId {
			isFollowed = true
			break
		}
	}

	// followedId的粉丝
	followerIDSli, err := repositories.QueryFollower(u.DB, followedId)
	if err != nil {
		return fmt.Errorf("UserReqContext.Follow() -> %v: ", err)
	}

	// 关注
	if !isFollowed {
		if err := repositories.InsertFollow(u.DB, followerId, followedId); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		//  followerId关注数量+1，followedId粉丝数量+1
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) + 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) + 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}

	} else { // 取消关注
		if _, err := repositories.DeleteObjectsByTable(u.DB, "user_follows", map[string]interface{}{"follower_id": followerId, "followed_id": followedId}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		//  followerId关注数量-1，followedId粉丝数量-1
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followerId}}, map[string]interface{}{"attention_count": len(followedIDSli) - 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
		if err = repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: followedId}}, map[string]interface{}{"fans_count": len(followerIDSli) - 1}); err != nil {
			return fmt.Errorf("UserReqContext.Follow() -> %v", err)
		}
	}

	return nil
}
