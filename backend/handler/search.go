package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/docs"
)

type SearchHandler struct {
	builder *docs.CachedBuilder
}

func NewSearchHandler(builder *docs.CachedBuilder) *SearchHandler {
	return &SearchHandler{builder: builder}
}

// GetSearch は GET /api/search のハンドラ。
// クエリパラメータ q で指定された文字列を含む記事を検索して返します。
func (h *SearchHandler) GetSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "検索キーワードが指定されていません"})
		return
	}

	tree, err := h.builder.BuildTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ツリーの構築に失敗しました"})
		return
	}

	results := docs.SearchTree(tree, query)
	if results == nil {
		results = []docs.SearchResult{}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
