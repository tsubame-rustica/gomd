package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
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

func initDB() {
	err := godotenv.Load()
	if err != nil {
		// ※Docker環境などで環境変数が直接渡される場合はエラーを無視してもOKな作りもあります
		log.Println(".envファイルが見つかりません。環境変数を使用します。")
	}

	// 🌟 os.Getenv で値を取得してDSNを組み立てる
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB接続エラー: ", err)
	}

	log.Println("データベース接続成功")
}

func main() {
	initDB()
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

	r.GET("/api/hello", func(c *gin.Context) {
		var jr PostResponse
		jr.Content = "Hello World!"
		c.JSON(200, jr)
	})

	// 静的ファイルの提供
	r.Static("/api/memo/content", "./content")

	r.GET("/api/memo/getPost/:category/:title", func(c *gin.Context) {
		category := c.Param("category")
		title := c.Param("title")

		mdPath := filepath.Join("content", category, title, title+".md")

		source, err := os.ReadFile(mdPath)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
			return
		}

		// MarkdownをHTMLに変換
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM), // テーブルや取り消し線などを有効化
		)

		var buf bytes.Buffer
		if err := md.Convert(source, &buf); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "変換エラー"})
			return
		}

		// 変換したHTMLをJSONで返す
		c.JSON(http.StatusOK, PostResponse{Content: buf.String()})
	})

	r.GET("/api/memo/getCategories", func(c *gin.Context) {
		// カテゴリ一覧を返す
		var categories []Category
		result := db.Find(&categories)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "データ取得エラー"})
			return
		}
		if len(categories) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "カテゴリが見つかりません"})
			return
		}

		var categoryList []CmnContentList

		for _, category := range categories {
			categoryList = append(categoryList, CmnContentList{
				Title: category.Name,
				Href:  fmt.Sprintf("/%s/default", category.DirName),
				Key:   fmt.Sprintf("%s", category.DirName),
			})
		}

		c.JSON(http.StatusOK, CmnResponse{Result: categoryList})
	})

	r.GET("/api/memo/getPosts/:category", func(c *gin.Context) {
		var posts []Post
		category := c.Param("category")

		var categoryId uint

		log.Println("category:", category)

		db.Model(&Category{}).Select("id").Where("dir_name = ?", category).Order("id asc").Limit(1).Scan(&categoryId)

		if categoryId == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "カテゴリが見つかりません"})
			return
		}
		log.Println("categoryId:", categoryId)

		result := db.Where("category_id = ?", categoryId).Find(&posts)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "データ取得エラー"})
			return
		}

		var postList []CmnContentList

		for _, post := range posts {
			filePath := fmt.Sprintf("/%s/%s", category, post.FileName)
			postList = append(postList, CmnContentList{
				Title: post.Title,
				Href:  filePath,
				Key:   post.FileName,
			})
		}

		log.Println("posts:", posts)
		c.JSON(http.StatusOK, CmnResponse{Result: postList})
	})

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}
