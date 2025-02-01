// Package utils 负责JWT的生成和验证
package utils

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 密钥
var jwtSecret = []byte("shadanchengzi") //用于生成和验证JWT的密钥

// Claims定义JWT的负载
type Claims struct {
	Username             string `json:"username"` //用户名
	jwt.RegisteredClaims        //包含JWT标准中定义的一些注册声明
}

// GenerateToken 生成JWT
func GenerateToken(username string) (string, error) { //根据传入的用户名生成一个JWT
	expirationTime := time.Now().Add(24 * time.Hour) //设置JWT的过期时间,24小时后过期
	claims := &Claims{                               //创建一个Claims结构体,包含用户名和过期时间
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{ //包含JWT标准中定义的一些注册声明
			ExpiresAt: jwt.NewNumericDate(expirationTime), //设置过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //创建一个包含Claims的JWT对象
	tokenString, err := token.SignedString(jwtSecret)          //使用密钥jwtSecret对JWT进行签名，生成最终的JWT字符串
	if err != nil {
		return "", err
	}
	return tokenString, nil //返回生成的JWT和可能的错误
}

// ValidateToken 验证JWT是否有效
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{} //创建一个空的Claims结构体,用于存储解析后的JWT负载
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil //返回密钥jwtSecret,用于验证JWT的签名
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid { //如果JWT无效,返回错误
		return nil, errors.New("invalid token")
	}
	return claims, nil //返回解析后的Claims和可能的错误
}

// Gin中间件， 用于处理请求时验证JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization") //从请求头中获取Authorization字段
		if authHeader == "" {                      //如果Authorization字段为空，返回401状态码和错误信息
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"}) //如果Authorization字段为空,返回401状态码和错误信息
			c.Abort()
			return
		}
		// 提取JWT
		token := strings.TrimPrefix(authHeader, "Bearer ") //从Authorization字段中提取JWT
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) //如果JWT为空，返回401状态码和错误信息
			c.Abort()
			return
		}
		// 验证JWT
		claims, err := ValidateToken(token) //验证JWT是否有效
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) //如果JWT无效，返回401状态码和错误信息
			c.Abort()
			return
		}
		// 将用户名添加到上下文中
		c.Set("username", claims.Username) //将用户名添加到上下文中
		c.Next()                           //继续处理下一个中间件
	}
}
