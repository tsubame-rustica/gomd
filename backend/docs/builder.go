package docs

import (
	"cmp"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const defaultOrder = math.MaxInt

// BuildTree は rootPath（contents/）配下の第一階層ディレクトリをカテゴリとして走査し、
// カテゴリ直下の .md ファイルのみを階層ツリーとして構築します（第二階層のサブディレクトリは含めません）。
func BuildTree(rootPath string) (*DocumentNode, error) {
	rootNode := &DocumentNode{
		Path:        rootPath,
		URLPath:     "/",
		DisplayName: filepath.Base(rootPath),
		Order:       defaultOrder,
		IsFile:      false,
		Children:    []*DocumentNode{},
	}

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		categoryPath := filepath.Join(rootPath, entry.Name())

		// 第一階層のカテゴリフォルダには _category.yml が必須
		meta, err := ParseCategoryYAML(categoryPath)
		if err != nil {
			continue
		}

		displayName := entry.Name()
		if meta.Category != "" {
			displayName = meta.Category
		}

		order := defaultOrder
		if meta.Order != 0 {
			order = meta.Order
		}

		categoryNode := &DocumentNode{
			Path:        categoryPath,
			URLPath:     "/" + filepath.ToSlash(entry.Name()),
			DisplayName: displayName,
			Order:       order,
			IsFile:      false,
			Children:    []*DocumentNode{},
		}

		// 第一階層ディレクトリ直下の .md ファイルのみを走査
		catEntries, err := os.ReadDir(categoryPath)
		if err != nil {
			return nil, err
		}

		for _, catEntry := range catEntries {
			// 第二階層ディレクトリはツリーに含めない
			if catEntry.IsDir() {
				continue
			}

			if strings.ToLower(filepath.Ext(catEntry.Name())) != ".md" {
				continue
			}

			filePath := filepath.Join(categoryPath, catEntry.Name())
			fmeta, err := ParseFrontmatter(filePath)
			if err != nil || strings.TrimSpace(fmeta.Title) == "" {
				// メタ情報（Frontmatter または title）がないファイルはスキップ
				continue
			}

			rel, _ := filepath.Rel(rootPath, filePath)
			fileURLPath := "/" + filepath.ToSlash(rel)

			fileOrder := defaultOrder
			if fmeta.Order != 0 {
				fileOrder = fmeta.Order
			}

			fileNode := &DocumentNode{
				Path:        filePath,
				URLPath:     fileURLPath,
				DisplayName: fmeta.Title,
				Order:       fileOrder,
				IsFile:      true,
				Children:    []*DocumentNode{},
			}

			categoryNode.Children = append(categoryNode.Children, fileNode)
		}

		// ファイルを Order 昇順でソート（cmp.Compare でオーバーフロー防止）
		slices.SortFunc(categoryNode.Children, func(a, b *DocumentNode) int {
			return cmp.Compare(a.Order, b.Order)
		})

		rootNode.Children = append(rootNode.Children, categoryNode)
	}

	// カテゴリを Order 昇順でソート（cmp.Compare でオーバーフロー防止）
	slices.SortFunc(rootNode.Children, func(a, b *DocumentNode) int {
		return cmp.Compare(a.Order, b.Order)
	})

	return rootNode, nil
}
