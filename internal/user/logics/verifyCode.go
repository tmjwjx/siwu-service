package logics

import (
	"errors"
	"forum/internal/internal_pkg/internal_utils"
	"gopkg.in/gomail.v2"
)

// VerifyCode 验证码结构体
type VerifyCode struct {
	from          string // 发送人
	to            string // 接收人
	subject       string // 主题
	body          string // 正文内容
	authorizeCode string // 授权码
}

func NewVerifyCode(from string, to string, subject string, body string, authorizeCode string) (*VerifyCode, error) {
	if internal_utils.IsValidEmail(from) {
		return nil, errors.New("发送人邮箱错误")
	}
	if internal_utils.IsValidEmail(to) {
		return nil, errors.New("接收人邮箱错误")
	}

	return &VerifyCode{
		from:          from,
		to:            to,
		subject:       subject,
		body:          body,
		authorizeCode: authorizeCode,
	}, nil
}

// SendEmail from 给 to 发送指定的邮件消息
func (vc *VerifyCode) SendEmail() error {
	m := gomail.NewMessage()
	// 设置邮件消息的头部字段
	m.SetHeader("From", vc.from)       // 发送人
	m.SetHeader("To", vc.to)           // 接收人
	m.SetHeader("Subject", vc.subject) // 主题
	m.SetBody("text/plain", vc.body)   // 正文内容
	// 创建一个新的邮件拨号器对象，用于通过指定的 SMTP 服务器发送邮件
	d := gomail.NewDialer("smtp.qq.com", 587, vc.from, vc.authorizeCode)
	// 通过拨号器对象发送指定的邮件消息
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
