package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用于存储用户数据的内存数据库
var db = make(map[string]string)

// setupRouter 初始化并配置路由:创建了三个主要的路由端点
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// 1.Ping测试路由，访问/ping时返回pong
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 2.获取用户值的路由，访问/user/:name时返回用户名和对应的值
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// 需要认证的路由组 (使用 gin.BasicAuth() 中间件)
	// 等同于:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	// 3.需要认证的路由组，访问/admin时需要基本认证
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // 用户名:foo 密码:bar
		"manu": "123", // 用户名:manu 密码:123
	}))

	/* 使用基本认证访问 /admin 的curl示例
	   Zm9vOmJhcg== 是 base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/

	// 管理员路由 - 需要基本认证
	authorized.POST("admin", func(c *gin.Context) {
		// 获取认证用户名
		user := c.MustGet(gin.AuthUserKey).(string)

		// 定义并解析JSON请求体
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		// 绑定JSON数据并存储到数据库
		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

// main 函数启动HTTP服务器
func main() {
	r := setupRouter()
	// 在 0.0.0.0:8080 监听并启动服务
	r.Run(":8080")
}
