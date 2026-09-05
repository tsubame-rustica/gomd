package docs

import (
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const defaultOrder = math.MaxInt

func BuildTree(rootPath string) (*DocumentNode, error) {
	return BuildNode(rootPath, rootPath)
}

func BuildNode(path, rootPath string) (*DocumentNode, error) {
	// URLPath: rootPath からの相対パスをスラッシュ区切りに変換
	rel, _ := filepath.Rel(rootPath, path)
	var urlPath string
	if rel == "." {
		urlPath = "/" // ルートノード
	} else {
		urlPath = "/" + filepath.ToSlash(rel)
	}

	node := &DocumentNode{
		Path:        path,
		URLPath:     urlPath,
		DisplayName: filepath.Base(path), // フォールバック: ディレクトリ名
		Order:       defaultOrder,
		IsFile:      false,
		Children:    []*DocumentNode{},
	}

	// _category.yml を読んで DisplayName と Order を上書き
	if meta, err := ParseCategoryYAML(path); err == nil {
		if meta.Category != "" {
			node.DisplayName = meta.Category
		}
		node.Order = meta.Order
	}

	// ディレクトリ内のエントリを走査
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// 第一階層（rootPath直下）のディレクトリは _category.yml が必須
			if path == rootPath {
				categoryPath := filepath.Join(entryPath, "_category.yml")
				if _, err := os.Stat(categoryPath); err != nil {
					// _category.yml が存在しない場合はスキップ
					continue
				}
			}

			// サブディレクトリ → 再帰的に BuildNode
			child, err := BuildNode(entryPath, rootPath)
			if err != nil {
				return nil, err
			}

			// 第二階層以降のディレクトリで、子要素（.mdなど）を持たない空フォルダ（画像専用フォルダなど）はスキップ
			if path != rootPath && len(child.Children) == 0 {
				continue
			}

			node.Children = append(node.Children, child)

		} else if strings.ToLower(filepath.Ext(entry.Name())) == ".md" {
			// Frontmatterをパース。メタ情報（Title）がないファイルはスキップ
			meta, err := ParseFrontmatter(entryPath)
			if err != nil || strings.TrimSpace(meta.Title) == "" {
				continue
			}

			rel, _ := filepath.Rel(rootPath, entryPath)
			fileURLPath := "/" + filepath.ToSlash(rel)

			fileNode := &DocumentNode{
				Path:        entryPath,
				URLPath:     fileURLPath,
				DisplayName: meta.Title,
				Order:       defaultOrder,
				IsFile:      true,
				Children:    []*DocumentNode{},
			}

			if meta.Order != 0 {
				fileNode.Order = meta.Order
			}

			node.Children = append(node.Children, fileNode)
		}
	}

	// Order 昇順でソート
	slices.SortFunc(node.Children, func(a, b *DocumentNode) int {
		return a.Order - b.Order
	})

	return node, nil
}
