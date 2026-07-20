package main

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type JsonResponse struct {
	Message string `json:"message"`
}

func main() {
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
		var jr JsonResponse
		jr.Message = "Hello World!"
		c.JSON(200, jr)
	})

	// 静的ファイルの提供
	r.Static("/api/memo/content", "./content")

	r.GET("api/memo/:category/:slug", func(c *gin.Context) {
		category := c.Param("category")
		slug := c.Param("slug")

		mdPath := filepath.Join("content", category, slug, slug+".md")

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
		c.JSON(http.StatusOK, JsonResponse{Message: buf.String()})
	})

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}
