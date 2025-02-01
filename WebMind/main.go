package main

import (
	"log"
	"webmind/handlers"
	"webmind/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库，实现代码在utils/db_utils.go中
	db, err := utils.ConnectPostgresql()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 创建默认的 gin 路由器
	r := gin.Default()

	// 定义登陆相关的路由
	logGroup := r.Group("/api")
	{
		logGroup.POST("/register", handlers.RegisterUser(db)) //注册
		logGroup.POST("/login", handlers.LoginUser(db))       //登录
		logGroup.POST("/logout", handlers.LogoutUser(db))     //登出
	}
	// 定义受保护的会话相关的路由
	sessionGroup := r.Group("/api/sessions")
	sessionGroup.Use(utils.AuthMiddleware())
	{
		// sessionGroup.POST("/", handlers.CreateSession(db)) //创建会话
		sessionGroup.POST("/", func(c *gin.Context) {
			// 获取存储在上下文中的用户名
			username, exists := c.Get("username")
			if !exists {
				c.JSON(500, gin.H{"error": "Username not found in context"})
				return
			}
			c.JSON(200, gin.H{"message": "JWT verification successful", "username": username})
		})
	}
	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	//关闭数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}
	defer sqlDB.Close()
}

/***
创建会话：
$body = @{
    web_pages = @("https://example.com", "https://another-example.com")
	} | ConvertTo-Json

	Invoke-WebRequest -Method Post -Uri "http://127.0.0.1:8080/sessions" -ContentType "application/json" -Body $body
获取会话：
Invoke-WebRequest -Method Get -Uri "http://127.0.0.1:8080/sessions/session_1"
简洁版：
$response = Invoke-WebRequest -Method Get -Uri "http://127.0.0.1:8080/sessions/session_1"
$response.Content

提交问题：
$questionBody = @{
    question = "WHAT is AI？"
} | ConvertTo-Json

Invoke-WebRequest -Method Post -Uri "http://127.0.0.1:8080/sessions/session_1/questions" -ContentType "application/json" -Body $questionBody
***/
