package token

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// 用于签名和验证 JWT 的密钥
var jwtSecret = []byte("siwu-web-service.forumSetJwtSecret_S@mpl3ComplexS3cretK3y")

// Claims 自定义的 Claims 结构体
type Claims struct {
	ID                   uint `json:"id"`
	jwt.RegisteredClaims      // 包含标准的 JWT 声明
}

// GenerateToken 使用用户的 ID 生成 JWT token。
func GenerateToken(id uint) (string, error) {
	// 创建声明 Claims
	claims := Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // 签发时间
			Issuer:    "siwu-web-service",                                // 签发者
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken 验证并解析 Token
func ValidateToken(tokenString string) (*Claims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的 token")
}

// AuthMiddleware JWT 认证中间件。在 JWT 验证通过后，将 ID 存储在上下文中，供后续路由使用。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 缺少授权标头 Authorization"), nil))
			// 中止剩余的中间件和处理函数执行，直接返回响应
			c.Abort()
			return
		}

		// 提取 Token 部分，去掉 "Bearer " 前缀
		tokenString = tokenString[len("Bearer "):]

		// 验证并解析 Token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 无效的 token"), nil))
			c.Abort()
			return
		}

		// 将用户 ID 保存到上下文中
		c.Set("id", claims.ID)
		c.Next()
	}
}
