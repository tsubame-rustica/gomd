package handler

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type ContentHandler struct {
	contentRoot string
}

func NewContentHandler(contentRoot string) *ContentHandler {
	return &ContentHandler{contentRoot: contentRoot}
}

// GetContent は GET /api/contents/*path のハンドラ。
// contents/ 配下の .md ファイルを読み込み、HTML に変換して返す。
func (h *ContentHandler) GetContent(c *gin.Context) {
	// *path パラメータを取得（例: /git/default/default.md）
	urlPath := c.Param("path")

	// パストラバーサル対策
	absRoot, err := filepath.Abs(h.contentRoot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー設定エラー"})
		return
	}

	requestedPath := filepath.Join(absRoot, filepath.FromSlash(urlPath))
	cleanPath := filepath.Clean(requestedPath)

	if !strings.HasPrefix(cleanPath, absRoot) {
		c.JSON(http.StatusForbidden, gin.H{"error": "アクセスが許可されていません"})
		return
	}

	ext := strings.ToLower(filepath.Ext(cleanPath))

	// 画像などの静的ファイルはそのまま返す
	staticExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".svg": true, ".webp": true, ".ico": true,
	}
	if staticExts[ext] {
		c.File(cleanPath)
		return
	}

	// .md 以外（ディレクトリ、不明な拡張子）はJSONエラーを返す
	if ext != ".md" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Markdown ファイルのみ取得できます"})
		return
	}

	source, err := os.ReadFile(cleanPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
		return
	}

	// Frontmatter（--- ... ---）を除去してからMarkdown変換する
	source = stripFrontmatter(source)

	// Markdown → HTML 変換（GFM + 脚注）
	md := goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Footnote))
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "変換エラー"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content":  buf.String(),
		"contents": buf.String(),
	})
}

// stripFrontmatter は --- で囲まれた Frontmatter ブロックを除去して残りのバイト列を返す。
func stripFrontmatter(src []byte) []byte {
	sep := []byte("---")
	newline := []byte("\n")

	// 1行目が --- でなければそのまま返す
	if !bytes.HasPrefix(src, sep) {
		return src
	}
	// 1行目の --- の後ろから次の --- を探す
	rest := src[len(sep):]
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	}
	idx := bytes.Index(rest, append(newline, sep...))
	if idx == -1 {
		return src
	}
	// 閉じ --- の行を飛ばした残りを返す
	after := rest[idx+len(newline)+len(sep):]
	return bytes.TrimLeft(after, "\n")
}

