package docs

import (
	"os"
	"strings"
)

// SearchTree はツリー全体を走査してキーワードにマッチする記事を返します
func SearchTree(node *DocumentNode, query string) []SearchResult {
	if query == "" {
		return nil
	}
	query = strings.ToLower(query)
	var results []SearchResult
	searchNode(node, query, &results)
	return results
}

func searchNode(node *DocumentNode, query string, results *[]SearchResult) {
	if node == nil {
		return
	}

	if node.IsFile {
		matched, snippet := checkMatch(node, query)
		if matched {
			*results = append(*results, SearchResult{
				URLPath:     node.URLPath,
				DisplayName: node.DisplayName,
				Snippet:     snippet,
			})
		}
	}

	for _, child := range node.Children {
		searchNode(child, query, results)
	}
}

func checkMatch(node *DocumentNode, query string) (bool, string) {
	// タイトルでマッチするかチェック
	if strings.Contains(strings.ToLower(node.DisplayName), query) {
		return true, "" // タイトルマッチの場合はスニペットなし（または適当なスニペット）
	}

	// 本文を読み込んでマッチするかチェック
	data, err := os.ReadFile(node.Path)
	if err != nil {
		return false, ""
	}
	
	content := string(data)
	contentLower := strings.ToLower(content)
	idx := strings.Index(contentLower, query)
	if idx != -1 {
		// スニペットを抽出 (前後30文字程度)
		start := idx - 30
		if start < 0 {
			start = 0
		}
		end := idx + len(query) + 30
		if end > len(content) {
			end = len(content)
		}
		
		snippet := content[start:end]
		// 改行をスペースに置換
		snippet = strings.ReplaceAll(snippet, "\n", " ")
		if start > 0 {
			snippet = "..." + snippet
		}
		if end < len(content) {
			snippet = snippet + "..."
		}
		return true, snippet
	}

	return false, ""
}
