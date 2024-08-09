package internal_utils

import (
	"regexp"
	"unicode"
)

// IsValidNickname 判断用户名是否合法
func IsValidNickname(nickname string) bool {
	// 用户名可以包含数字、字母及中文，长度不超过 16 个字符
	re := regexp.MustCompile(`^[\u4e00-\u9fa5A-Za-z0-9]{1,16}$`)
	return re.MatchString(nickname)
}

// IsValidEmail 判断邮箱是否合法。
func IsValidEmail(email string) bool {
	// 定义正则表达式
	qqEmailPattern := `^[1-9][0-9]{4,10}@qq\.com$`
	// 编译正则表达式
	re := regexp.MustCompile(qqEmailPattern)

	return re.MatchString(email)
}

// IsValidPassword 判断密码是否合法。密码只能也必须同时包含数字和字母，长度在6到16位之间。
func IsValidPassword(password string) bool {
	pattern := "^[a-zA-Z0-9]{6,16}$"
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(password) {
		return false
	}

	// 检查密码是否包含至少一个字母和一个数字
	hasLetter := false
	hasDigit := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
		// 提前退出
		if hasLetter && hasDigit {
			return true
		}
	}
	return hasLetter && hasDigit
}
