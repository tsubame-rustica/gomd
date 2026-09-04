package docs_test

import (
	"os"
	"path/filepath"
	"testing"

	"backend/docs"
)

// テスト用のコンテンツディレクトリを一時的に作成するヘルパー
func setupTestDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// root/
	// ├── _category.yml    (label: Root, order: 1)
	// ├── git/
	// │   ├── _category.yml (label: Git, order: 10)
	// │   ├── intro.md      (title: はじめに, order: 5)
	// │   └── branch.md     (title: ブランチ, order: 20)
	// └── linux/
	//     ├── _category.yml (label: Linux, order: 20)
	//     └── basics.md     (title: 基本コマンド, order: 10)

	mustWriteFile(t, filepath.Join(root, "_category.yml"), "label: \"Root\"\norder: 1\n")

	gitDir := filepath.Join(root, "git")
	mustMkdir(t, gitDir)
	mustWriteFile(t, filepath.Join(gitDir, "_category.yml"), "label: \"Git\"\norder: 10\n")
	mustWriteFile(t, filepath.Join(gitDir, "intro.md"), "---\ntitle: \"はじめに\"\norder: 5\n---\n\n# Gitについて\n")
	mustWriteFile(t, filepath.Join(gitDir, "branch.md"), "---\ntitle: \"ブランチ\"\norder: 20\n---\n\n# ブランチ\n")

	linuxDir := filepath.Join(root, "linux")
	mustMkdir(t, linuxDir)
	mustWriteFile(t, filepath.Join(linuxDir, "_category.yml"), "label: \"Linux\"\norder: 20\n")
	mustWriteFile(t, filepath.Join(linuxDir, "basics.md"), "---\ntitle: \"基本コマンド\"\norder: 10\n---\n\n# コマンド\n")

	return root
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// BuildTree がルートノードを返すこと
func TestBuildTree_ReturnsRoot(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}
	if node == nil {
		t.Fatal("BuildTree returned nil")
	}
}

// ルートノードの URLPath が "/" になること
func TestBuildTree_RootURLPath(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}
	if node.URLPath != "/" {
		t.Errorf("root URLPath = %q, want %q", node.URLPath, "/")
	}
}

// _category.yml の label が DisplayName に反映されること
func TestBuildTree_CategoryDisplayName(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}

	// git カテゴリを探す
	var gitNode *docs.DocumentNode
	for _, child := range node.Children {
		if !child.IsFile && child.DisplayName == "Git" {
			gitNode = child
			break
		}
	}
	if gitNode == nil {
		t.Fatal("git カテゴリが見つからない")
	}
}

// カテゴリが order 昇順にソートされていること
func TestBuildTree_CategoryOrder(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}

	dirs := []*docs.DocumentNode{}
	for _, child := range node.Children {
		if !child.IsFile {
			dirs = append(dirs, child)
		}
	}
	if len(dirs) < 2 {
		t.Fatalf("サブディレクトリが2つ以上必要: got %d", len(dirs))
	}
	// git(order:10) → linux(order:20) の順
	if dirs[0].DisplayName != "Git" || dirs[1].DisplayName != "Linux" {
		t.Errorf("order ソートが正しくない: got [%s, %s], want [Git, Linux]", dirs[0].DisplayName, dirs[1].DisplayName)
	}
}

// Frontmatter の title と order が .md ノードに反映されること
func TestBuildNode_FileMetadata(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}

	var gitNode *docs.DocumentNode
	for _, child := range node.Children {
		if !child.IsFile && child.DisplayName == "Git" {
			gitNode = child
			break
		}
	}
	if gitNode == nil {
		t.Fatal("git カテゴリが見つからない")
	}

	// intro.md (order:5) が先、branch.md (order:20) が後
	if len(gitNode.Children) < 2 {
		t.Fatalf("git 配下のファイルが2つ必要: got %d", len(gitNode.Children))
	}
	first := gitNode.Children[0]
	if first.DisplayName != "はじめに" {
		t.Errorf("最初のファイル DisplayName = %q, want \"はじめに\"", first.DisplayName)
	}
	if first.Order != 5 {
		t.Errorf("最初のファイル Order = %d, want 5", first.Order)
	}
	if !first.IsFile {
		t.Error("IsFile = false, want true")
	}
}

// IsFile フラグが正しく設定されること
func TestBuildNode_IsFile(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}

	for _, child := range node.Children {
		if child.IsFile {
			t.Errorf("root 直下に IsFile=true のノードがあった: %s", child.URLPath)
		}
	}
}

// URLPath がスラッシュ始まりになること
func TestBuildNode_URLPath(t *testing.T) {
	root := setupTestDir(t)
	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}
	for _, cat := range node.Children {
		for _, file := range cat.Children {
			if !file.IsFile {
				continue
			}
			if len(file.URLPath) == 0 || file.URLPath[0] != '/' {
				t.Errorf("URLPath %q がスラッシュ始まりでない", file.URLPath)
			}
		}
	}
}

// .md 以外のファイルはツリーに含まれないこと
func TestBuildTree_IgnoresNonMdFiles(t *testing.T) {
	root := setupTestDir(t)
	gitDir := filepath.Join(root, "git")
	mustWriteFile(t, filepath.Join(gitDir, "image.png"), "dummy")

	node, err := docs.BuildTree(root)
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}
	var gitNode *docs.DocumentNode
	for _, child := range node.Children {
		if !child.IsFile && child.DisplayName == "Git" {
			gitNode = child
			break
		}
	}
	for _, child := range gitNode.Children {
		if child.IsFile && filepath.Ext(child.URLPath) == ".png" {
			t.Errorf(".png ファイルがツリーに含まれている: %s", child.URLPath)
		}
	}
}
