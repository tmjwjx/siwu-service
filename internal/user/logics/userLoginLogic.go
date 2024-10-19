package logics

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strings"
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

	// 更新用户密码
	// err = repositories.UpdateObjects(u.DB, &models.User{}, map[string]interface{}{"id": userID, "email": email}, map[string]interface{}{"password": encryptedPassword})
	err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: userID}, Email: email}, map[string]interface{}{"password": encryptedPassword})
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}

	// 删除该用户对应的全部验证码
	_, err = sqlUtils.DeleteObjectsByModel(u.DB, &models.UserVerifyCode{}, map[string]interface{}{"user_id": userID})
	if err != nil {
		return fmt.Errorf("UserReqContext.Register() err: %v", err)
	}
	return nil
}

// ReqVerifyCode 用户请求验证码
func (u *UserReqContext) ReqVerifyCode(email string) error {
	// 随机生成验证码
	verifyCode := internalUtils.RandomGenerateStrings(internalUtils.VerifyCodeLen)

	// 存储数据

	// 判断该email是否已经有用户使用过
	user := repositories.QueryUserByEmail(u.DB, email)
	// 如果没有用户使用过这个email，向user表中插入用户，并向UserVerifyCode表中插入验证码
	if user == nil {
		// 随机生成用户名
		name := internalUtils.RandomGenerateStrings(internalUtils.UserNameLen)
		// 给用户生成一个默认密码
		password := internalUtils.RandomGenerateStrings(12)
		encryptedPassword, err := internalUtils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() : 密码%s加密失败", password)
		}
		lastLogintime := time.Now()
		// 使用InsertObject()方法向user表中插入新数据，model参数必须是指针类型
		if err := sqlUtils.InsertObject(u.DB, &models.User{Nickname: name, Email: email, Password: encryptedPassword, LastLoginTime: lastLogintime}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
		}

		// 查询该email对应的id
		us := repositories.QueryUserByEmail(u.DB, email)
		// 向 UserDetail 用户详情表中插入数据
		if err := sqlUtils.InsertObject(u.DB, &models.UserDetail{UserID: us.ID}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
		}

		// 获取用户id
		user = repositories.QueryUserByEmail(u.DB, email)
		// 将验证码插入到 UserVerifyCode表
		if err := sqlUtils.InsertObject(u.DB, &models.UserVerifyCode{UserID: user.ID, VerifyCode: verifyCode}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
		}

	} else { // 如果已经有用户使用，而且发送验证码的冷却时间到了，插入一条数据
		userVerifyCode, err := repositories.QueryLastUserVerifyCodeByUserID(u.DB, user.ID)
		if err != nil { // 执行错误，没有查询到验证码（可能是手动删除了数据库中的验证码，所以报错）
			// 插入一条新的验证码数据
			if err = sqlUtils.InsertObject(u.DB, &models.UserVerifyCode{UserID: user.ID, VerifyCode: verifyCode}); err != nil {
				return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
			}
			// 给用户发送验证码
			// body := fmt.Sprintf("你的验证码为 %s，有效时间为 %d 分钟\n", verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
			// 读取邮件模板
			templateFile, err := os.Open("internal/internalPkg/template/emailFormatTemplate.html")
			if err != nil {
				return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 无法打开模板文件: %v", err)
			}
			defer templateFile.Close()
			templateContent, err := ioutil.ReadAll(templateFile)
			if err != nil {
				return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 无法读取模板内容: %v", err)
			}
			// 格式化邮件内容
			body := fmt.Sprintf(string(templateContent), verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
			if err = u.SendEmail(email, internalUtils.VerifyCodeSubject, body); err != nil {
				return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
			}
			return nil
		}

		// 判断冷却时间
		// 当前时间
		now := time.Now()
		// 计算更新时间和当前时间的差异
		duration := now.Sub(userVerifyCode.UpdatedAt)
		// 如果冷却时间未到，返回错误
		if duration < internalUtils.VerifyCodeCoolTime {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 发送验证码正在冷却时间中")
		}

		// 插入一条新的验证码数据
		if err = sqlUtils.InsertObject(u.DB, &models.UserVerifyCode{UserID: user.ID, VerifyCode: verifyCode}); err != nil {
			return fmt.Errorf("UserReqContext.VerifyCodeReq() -> %v", err)
		}
	}

	// 给用户发送验证码
	// body := fmt.Sprintf("你的验证码为 %s，不区分大小写，有效时间为 %d 分钟\n", verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
	// 读取邮件模板
	templateFile, err := os.Open("internal/internalPkg/template/emailFormatTemplate.html")
	if err != nil {

		// test 目录位置
		currentDirectory, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return fmt.Errorf("Current Directory:", currentDirectory)

		//return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 无法打开模板文件: %v", err)

	}
	defer templateFile.Close()

	templateContent, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return fmt.Errorf("UserReqContext.VerifyCodeReq() err: 无法读取模板内容: %v", err)
	}
	// 格式化邮件内容
	body := fmt.Sprintf(string(templateContent), verifyCode, int(internalUtils.VerifyCodeEffectiveDuration.Minutes()))
	if err = u.SendEmail(email, internalUtils.VerifyCodeSubject, body); err != nil {
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
func (u *UserReqContext) Login(logicMsg requests.LogicReq) error {
	// 判断邮箱和密码是否匹配
	email := logicMsg.Email
	password := logicMsg.Password

	// 根据邮箱查用户
	user := repositories.QueryUserByEmail(u.DB, email)
	if user == nil {
		return fmt.Errorf("UserReqContext.Login() err: 不存在该邮箱用户")
	}

	// 比较加密密码
	encryptedPassword := user.Password
	if !internalUtils.CheckPasswordHash(password, encryptedPassword) {
		return fmt.Errorf("UserReqContext.Login() err: 密码错误")
	}

	// 改变 LastLoginTime
	now := time.Now() // 获取当前时间
	if err := sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: user.ID}}, map[string]interface{}{"last_login_time": now}); err != nil {
		return fmt.Errorf("UserReqContext.Login() -> %v", err)
	}

	return nil
}