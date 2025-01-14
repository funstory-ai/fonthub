package main

import (
	"github.com/gin-gonic/gin"
)

// gin server
func ginServer() {
	// Initialize gin router
	r := gin.Default()

	// 添加静态文件服务
	r.Static("/static", "./static")

	// 修改路由返回HTML而不是JSON
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Hello Fonts",
		})
	})

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// Run the server on port 8080
	r.Run(":8080")
}

func main() {
	ginServer()
}
