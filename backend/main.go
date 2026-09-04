package main

import (
	"backend/docs"
	"backend/handler"

	"github.com/gin-gonic/gin"
)

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

	// ドキュメントツリーのキャッシュビルダーを初期化
	builder := docs.NewCachedBuilder("./content")
	treeHandler := handler.NewTreeHandler(builder)
	contentHandler := handler.NewContentHandler("./content")
	searchHandler := handler.NewSearchHandler(builder)

	r.GET("/api/tree", treeHandler.GetTree)
	r.GET("/api/content/*path", contentHandler.GetContent)
	r.GET("/api/memo/content/*path", contentHandler.GetContent) // 過去のマークダウン互換用エイリアス
	r.GET("/api/search", searchHandler.GetSearch)

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}
