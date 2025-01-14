package main

import (
	"github.com/funstory-ai/fonthub/internal/handlers"
	"github.com/gin-gonic/gin"
)

// gin server
func ginServer() {
	// Initialize gin router
	r := gin.Default()

	// 加载HTML模板 (移到路由设置之前)
	r.LoadHTMLGlob("templates/*")

	// 添加静态文件服务
	r.Static("/static", "./static")

	// 基础路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Hello Fonts",
		})
	})

	// Create API route group
	api := r.Group("/api")
	{
		// 整合所有API路由
		api.GET("/fonts", handlers.GetFontsHandler)
		api.GET("/fonts/selector", handlers.GetFontsBySelectorHandler)
	}

	// Run the server on port 8080
	r.Run(":8080")
}

func main() {
	// 只需要调用 ginServer
	ginServer()
}
