package logic

import (
	"encoding/json"
	"fmt"
	"forum/internal/user/requests"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserLogic 登陆模块，包含注册，登陆，验证码
type UserLogic struct {
	DB  *gorm.DB
	Ctx *gin.Context
}

func NewUserLogic(db *gorm.DB, c *gin.Context) *UserLogic {
	return &UserLogic{
		DB:  db,
		Ctx: c,
	}
}

// Register 注册
func (u *UserLogic) Register() *response.AppErr {
	// 获取数据包
	data, err := u.Ctx.GetRawData()
	// 请求语法错误或无效参数
	if err != nil {
		return response.StatusBadRequestErr
	}

	fmt.Println("Register 接收到的消息为:")
	fmt.Println(string(data))

	// 将数据包反序列化
	var registerMsg requests.RegisterMsg
	err = json.Unmarshal(data, &registerMsg)
	// 反序列化失败
	if err != nil {
		return response.StatusInternalServerErr
	}

	// 验证码核对，核对两次输入的密码
	// TODO

	// 添加到数据库中...

	return nil
}
