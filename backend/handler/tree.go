package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/docs"
)

type TreeHandler struct {
	builder *docs.CachedBuilder
}

func NewTreeHandler(builder *docs.CachedBuilder) *TreeHandler {
	return &TreeHandler{builder: builder}
}

// GetTree は GET /api/tree のハンドラ。
// ドキュメントツリーを JSON で返す。
func (h *TreeHandler) GetTree(c *gin.Context) {
	tree, err := h.builder.BuildTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ツリーの構築に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, tree)
}
