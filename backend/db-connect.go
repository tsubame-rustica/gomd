package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() {
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
