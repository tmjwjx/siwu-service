package utils

import (
	"forum/internal/user/requests"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func UserDateVerify(user *requests.User) error {
	// 初始化验证器
	validate := validator.New()
	err := validate.RegisterValidation("nickname", NicknameValidation)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("email", EmailValidation)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("password", PasswordValidation)
	if err != nil {
		return err
	}
	// 参数验证
	if err := validate.Struct(user); err != nil {
		return err
	}
	return nil
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
	// 密码只能包含密码和字母，数字和字母至少各包含一个，并且长度不少于6，不大于16
	re := regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,16}$`)
	return re.MatchString(password)
}
