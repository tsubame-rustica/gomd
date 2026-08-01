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
)

func GetPost(c *gin.Context) {
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

	c.JSON(http.StatusOK, PostResponse{Content: buf.String()})
}

func GetCategories(c *gin.Context) {
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
}

func GetPosts(c *gin.Context) {
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
}
