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
			// サブディレクトリ → 再帰的に BuildNode
			child, err := BuildNode(entryPath, rootPath)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)

		} else if strings.ToLower(filepath.Ext(entry.Name())) == ".md" {
			// .md ファイル → Frontmatter を読んでノード生成
			rel, _ := filepath.Rel(rootPath, entryPath)
			fileURLPath := "/" + filepath.ToSlash(rel)

			fileNode := &DocumentNode{
				Path:        entryPath,
				URLPath:     fileURLPath,
				DisplayName: strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())), // フォールバック
				Order:       defaultOrder,
				IsFile:      true,
				Children:    []*DocumentNode{},
			}

			if meta, err := ParseFrontmatter(entryPath); err == nil {
				if meta.Title != "" {
					fileNode.DisplayName = meta.Title
				}
				if meta.Order != 0 {
					fileNode.Order = meta.Order
				}
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
