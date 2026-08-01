package main

import (
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"backend/admin"
)

var db *gorm.DB

type PostResponse struct {
	Content string `json:"content"`
}

type CmnContentList struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Key   string `json:"key"`
}

type CmnResponse struct {
	Result []CmnContentList `json:"result"`
}

type Category struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	DirName string `json:"dir_name"`
}

type Post struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Title      string `json:"title"`
	FileName   string `json:"file_name"`
	CategoryID uint   `json:"category_id"`
}

func main() {
	InitDB()
	r := gin.Default()

	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// ブラウザからの事前確認(OPTIONSリクエスト)にはすぐOK(204)を返す
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	adminGroup := r.Group("/api/admin")
	{
		adminGroup.GET("/login")

		// カテゴリの作成・更新・削除 (CRUD)
		adminGroup.POST("/categories", admin.CreateCategory)
		adminGroup.PUT("/categories/:id", admin.UpdateCategory)
		adminGroup.DELETE("/categories/:id", admin.DeleteCategory)

		// 記事の作成・更新・削除
		adminGroup.POST("/posts", admin.CreatePost)
		adminGroup.PUT("/posts/:id", admin.UpdatePost)
		adminGroup.DELETE("/posts/:id", admin.DeletePost)
	}

	r.GET("/api/hello", func(c *gin.Context) {
		var jr PostResponse
		jr.Content = "Hello World!"
		c.JSON(200, jr)
	})

	// 静的ファイルの提供
	r.Static("/api/content", "./content")

	r.GET("/api/getPost/:category/:title", GetPost)

	r.GET("/api/getCategories", GetCategories)

	r.GET("/api/getPosts/:category", GetPosts)

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}
