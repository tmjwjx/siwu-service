package controllers

import (
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
)

func PersonalDataHandler(c *gin.Context) {
	// 初始化验证器
	validate := validator.New()
	err := validate.RegisterValidation("nickname", NicknameValidation)
	if err != nil {
		return
	}
	err = validate.RegisterValidation("email", EmailValidation)
	if err != nil {
		return
	}
	err = validate.RegisterValidation("password", PasswordValidation)
	if err != nil {
		return
	}

	// 将前端传来的请求数据绑定到 User 结构体
	var user requests.User
	if err := c.ShouldBind(&user); err != nil {
		// 处理绑定错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 参数验证
	if err := validate.Struct(user); err != nil {
		// 输出错误信息
		//fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "用户名或密码的格式不正确"})
		return
	}

	// 业务处理
	logics.PersonalDataLogic(&user, c)

}

// 自定义的验证函数

// NicknameValidation Nickname 验证函数
func NicknameValidation(fl validator.FieldLevel) bool {
	nickname := fl.Field().String()
	// 用户名可以包含数字、字母及中文，长度不超过 16 个字符
	re := regexp.MustCompile(`^[\u4e00-\u9fa5A-Za-z0-9]{1,16}$`)
	return re.MatchString(nickname)
}

// EmailValidation Email 验证函数
func EmailValidation(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	// 检查电子邮件的格式是否正确
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// PasswordValidation Password 验证函数
func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// 密码长度不超过 16 个字符，只能包含数字和字母
	re := regexp.MustCompile(`^[A-Za-z0-9]{1,16}$`)
	return re.MatchString(password)
}
