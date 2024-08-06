package logics

import (
	"errors"
	"forum/internal/internal_pkg/internal_utils"
	"gopkg.in/gomail.v2"
)

// SendEmail from 给 to 发送指定的邮件消息
func SendEmail(from string, to string, subject string, body string, authorizeCode string) error {
	if internal_utils.IsValidEmail(from) {
		return errors.New("发送人邮箱错误")
	}
	if internal_utils.IsValidEmail(to) {
		return errors.New("接收人邮箱错误")
	}

	m := gomail.NewMessage()
	// 设置邮件消息的头部字段
	m.SetHeader("From", from)       // 发送人
	m.SetHeader("To", to)           // 接收人
	m.SetHeader("Subject", subject) // 主题
	m.SetBody("text/plain", body)   // 正文内容
	// 创建一个新的邮件拨号器对象，用于通过指定的 SMTP 服务器发送邮件
	d := gomail.NewDialer("smtp.qq.com", 587, from, authorizeCode)
	// 通过拨号器对象发送指定的邮件消息
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
