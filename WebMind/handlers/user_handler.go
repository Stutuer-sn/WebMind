// 负责用户的注册、登录、登出、获取用户信息等操作
package handlers

import (
	"encoding/json"
	"net/http"
	"webmind/models"
	"webmind/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUser 处理用户注册请求
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		//	从 HTTP 请求体中解析 JSON 数据,并映射到models.User 结构体实例 user 上
		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//	对用户输入的明文密码进行哈希处理
		if err := user.SetPassword(user.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to set password"})
			return
		}
		//	检查用户名是否已存在
		if err := db.Where("username = ?", user.Username).First(&models.User{}).Error; err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		//使用 GORM 的 Create 方法将用户信息插入到数据库中
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfull"})
	}
}

// 处理用户登录请求
func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 从数据库中读取用户信息
		var storedUser models.User
		if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		if err := storedUser.CheckPassword(user.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

		token, err := utils.GenerateToken(storedUser.Username) //生成JWT
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// 处理用户登出请求
func LogoutUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
