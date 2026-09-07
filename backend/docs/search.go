package docs

import (
	"strings"
)

// SearchTree はツリー全体を走査してキーワードにマッチする記事を返します。
// contents はファイルパス→Markdown本文のキャッシュマップで、ディスクI/Oを回避します。
func SearchTree(node *DocumentNode, query string, contents map[string]string) []SearchResult {
	if query == "" {
		return nil
	}
	query = strings.ToLower(query)
	var results []SearchResult
	searchNode(node, query, contents, &results)
	return results
}

func searchNode(node *DocumentNode, query string, contents map[string]string, results *[]SearchResult) {
	if node == nil {
		return
	}

	if node.IsFile {
		matched, snippet := checkMatch(node, query, contents)
		if matched {
			*results = append(*results, SearchResult{
				URLPath:     node.URLPath,
				DisplayName: node.DisplayName,
				Snippet:     snippet,
			})
		}
	}

	for _, child := range node.Children {
		searchNode(child, query, contents, results)
	}
}

func checkMatch(node *DocumentNode, query string, contents map[string]string) (bool, string) {
	// タイトルでマッチするかチェック
	if strings.Contains(strings.ToLower(node.DisplayName), query) {
		return true, "" // タイトルマッチの場合はスニペットなし
	}

	// キャッシュから本文を取得（ディスクI/Oなし）
	content, ok := contents[node.Path]
	if !ok {
		return false, ""
	}

	contentLower := strings.ToLower(content)
	byteIdx := strings.Index(contentLower, query)
	if byteIdx == -1 {
		return false, ""
	}

	// UTF-8 安全なスニペット抽出:
	// byteIdx はバイト位置のため、そのままスライスすると日本語などマルチバイト文字の
	// 途中でスライスして不正な文字列になる恐れがある。rune 変換して文字境界で切り出す。
	runes := []rune(content)
	runeIdx := len([]rune(content[:byteIdx]))
	queryRuneLen := len([]rune(query))

	const context = 30 // 前後に表示する文字数
	start := runeIdx - context
	if start < 0 {
		start = 0
	}
	end := runeIdx + queryRuneLen + context
	if end > len(runes) {
		end = len(runes)
	}

	snippet := string(runes[start:end])
	snippet = strings.ReplaceAll(snippet, "\n", " ")
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(runes) {
		snippet += "..."
	}
	return true, snippet
}
