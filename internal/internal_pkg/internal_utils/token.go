package internal_utils

import (
	"fmt"
	"forum/internal/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 用于签名和验证 JWT 的密钥
var jwtSecret = []byte("siwu-web-service.forumSetJwtSecret")

// GenerateToken 根据用户的用户名和密码产生token
func GenerateToken(email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)
	user := models.User{
		Email:    email,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "liuqi-forum-app",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return token, fmt.Errorf("GenerateToken() err: %v", err)
	}
	return token, nil
}

// ParseToken 根据传入的token值获取到Claims对象信息（进而获取其中的用户名和密码）
func ParseToken(token string) (*models.User, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &models.User{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if user, ok := tokenClaims.Claims.(*models.User); ok && tokenClaims.Valid {
			return user, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("ParseToken() err: %v", err)
	}

	return nil, nil
}
