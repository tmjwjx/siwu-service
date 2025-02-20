package sendEmailAsynchronous

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"forum/pkg/response"
	"gopkg.in/gomail.v2"
)

// SendEmail 给邮箱(to)发送内容(body)
// to: 接收人
// body: 正文内容
// subject: 主题
func SendEmail(s *globals.SendEmailConfig, to string, subject string, body string) error {
	// 判断邮箱是否合法
	if !internalUtils.IsValidEmail(to) {
		globals.Log.Errorf(response.ErrEmailIsInvalid + ":" + to)
		return fmt.Errorf(response.ErrEmailIsInvalid + ":" + to)
	}

	m := gomail.NewMessage()
	// 设置邮件消息的头部字段
	m.SetHeader("From", s.From)     // 发送人
	m.SetHeader("To", to)           // 接收人
	m.SetHeader("Subject", subject) // 主题
	m.SetBody("text/html", body)    // 正文内容
	// 创建一个新的邮件拨号器对象，用于通过指定的 SMTP 服务器发送邮件
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.AuthorizeCode)
	// 通过拨号器对象发送指定的邮件消息
	if err := d.DialAndSend(m); err != nil {
		// return fmt.Errorf("UserReqContext.SendEmail() err: %v", err)
		globals.Log.Errorf(err.Error())
		return err
	}
	return nil
}
