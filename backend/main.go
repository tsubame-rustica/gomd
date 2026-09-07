package main

import (
	"log"
	"os"

	"backend/docs"
	"backend/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS設定（環境変数 CORS_ALLOW_ORIGIN があれば指定）
	allowedOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
	if allowedOrigin == "" {
		if gin.Mode() == gin.DebugMode {
			// 開発時のみ: Vite のデフォルトポートへのフォールバック
			allowedOrigin = "http://localhost:5173"
			log.Println("[WARN] CORS_ALLOW_ORIGIN が未設定です。デバッグモードのため http://localhost:5173 を使用します。")
		} else {
			// 本番環境では CORS_ALLOW_ORIGIN の設定必須
			log.Fatal("[ERROR] 本番環境では CORS_ALLOW_ORIGIN を必ず設定してください。")
		}
	}

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// ブラウザからの事前確認(OPTIONSリクエスト)にはすぐOK(204)を返す
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ドキュメントツリーのキャッシュビルダーを初期化
	builder := docs.NewCachedBuilder("./contents")
	treeHandler := handler.NewTreeHandler(builder)
	contentHandler := handler.NewContentHandler("./contents")
	searchHandler := handler.NewSearchHandler(builder)

	r.GET("/api/tree", treeHandler.GetTree)
	r.GET("/api/contents/*path", contentHandler.GetContent)
	r.GET("/api/search", searchHandler.GetSearch)

	// Cloud Run 等の環境変数 PORT に対応（デフォルト: 8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
