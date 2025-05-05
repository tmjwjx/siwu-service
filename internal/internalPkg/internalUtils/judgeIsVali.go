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
	// 定义通用的邮箱正则表达式
	// 这个正则表达式可以匹配大多数常见的邮箱格式
	// 它允许用户名部分包含字母、数字、点、加号、减号和下划线
	// 域名部分允许字母、数字、连字符，并且以有效的顶级域名（如.com、.org 等）结尾
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// 编译正则表达式
	re := regexp.MustCompile(emailPattern)

	return re.MatchString(email)
}

// IsValidPassword 判断密码是否合法。密码必须包含英文字母、数字、特殊字符，长度在8到20位之间，不限制中文字符。
func IsValidPassword(password string) bool {
	// 密码长度在8到20位之间
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	hasDigit := false         // 英文字母
	hasSpecialChar := false   // 数字
	hasEnglishLetter := false // 符号（运算符、特殊字符）
	// 遍历密码中的每个字符
	for _, char := range password {
		if unicode.Is(unicode.Latin, char) { // 确保是英文字符
			hasEnglishLetter = true
		} else if unicode.IsDigit(char) { // 确保是数字
			hasDigit = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) { // unicode.IsPunct：是否是标点符号。unicode.IsSymbol：是否是符号。
			hasSpecialChar = true
		}
	}
	return hasEnglishLetter && hasDigit && hasSpecialChar
}
