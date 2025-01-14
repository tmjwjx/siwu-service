package internalUtils

import (
	"regexp"
	"unicode"
)

// IsValidNickname 判断用户名是否合法
func IsValidNickname(nickname string) bool {
	// 用户名可以包含数字、字母及中文，长度不超过 16 个字符
	re := regexp.MustCompile(`^[\p{Han}A-Za-z0-9]{1,16}$`)
	return re.MatchString(nickname)
}

// IsValidEmail 判断邮箱是否合法。
func IsValidEmail(email string) bool {
	// 定义正则表达式，允许字母、数字和下划线
	qqEmailPattern := `^[a-zA-Z0-9_]{5,15}@qq\.com$`
	// 编译正则表达式
	re := regexp.MustCompile(qqEmailPattern)

	return re.MatchString(email)
}

// IsValidPassword 判断密码是否合法。密码只能也必须同时包含数字和字母，长度在8到16位之间。
func IsValidPassword(password string) bool {
	pattern := "^[a-zA-Z0-9!@#$%^&*()_+\\-={}|\\[\\]:\";'<>?,./]{8,16}$"
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(password) {
		return false
	}

	// 密码必须要包含字母、数字、特殊字符
	hasLetter := false
	hasDigit := false
	hasSpecialChar := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialChar = true
		}
		// 提前退出
		if hasLetter && hasDigit && hasSpecialChar {
			return true
		}
	}
	return hasLetter && hasDigit && hasSpecialChar
}
