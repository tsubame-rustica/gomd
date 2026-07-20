package main

import (
	"github.com/gin-gonic/gin"
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

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}
