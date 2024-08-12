package internal_utils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

// RandomGenerateStrings 随机生成长度为 l 的数字字母混合的字符串
func RandomGenerateStrings(l int) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// HashPassword 加密密码。使用 bcrypt.GenerateFromPassword() 方法来加密密码。
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash 验证密码。使用 bcrypt.CompareHashAndPassword() 方法来验证输入的密码是否正确。
// password: 要比较的密码。
// hashedPassword: 原先的哈希密码。
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
